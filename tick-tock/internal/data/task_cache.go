package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"tick-tock/internal/biz"
	"tick-tock/internal/conf"
	"tick-tock/pkg/log"
	"tick-tock/util/task"
	"time"
)

type taskCache struct {
	confData *conf.Data
	data     *Data
}

func (c *taskCache) GetValByStartAndEnd(ctx context.Context, key string, startUnixMilli int64, endUnixMilli int64) ([]*biz.Task, error) {
	results, err := c.data.redis.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", startUnixMilli),
		Max: fmt.Sprintf("%d", endUnixMilli),
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Warn(ctx, "key is not found.", "error", err, "key", key)
			return nil, err
		}
		log.Error(ctx, "get val from redis error.", "error", err, "key", key)
		return nil, err
	}
	length := len(results)
	tasks := make([]*biz.Task, length)
	for i := range results {
		tID, runTime, err := task.SplitTimerIDAndRunTime(results[i])
		if err != nil {
			log.Error(ctx, "split timer id and run time error.", "error", err)
			return nil, err
		}
		tasks[i] = &biz.Task{
			Tid:     tID,
			RunTime: time.UnixMilli(runTime),
		}
	}
	log.Info(ctx, "get val from redis success.", "key", key, "count", length)
	return tasks, nil
}

func (c *taskCache) SaveTasks(ctx context.Context, tasks []*biz.Task) error {
	pipe := c.data.redis.Pipeline()
	if _, err := pipe.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := range tasks {
			// 不同索引对应的运行时间不同
			// tableName 也不同，不能进行聚合写优化
			tableName := task.GetTableName(ctx, c.confData, tasks[i].RunTime, tasks[i].ID)
			pipe.ZAdd(ctx, tableName, redis.Z{
				Score:  float64(tasks[i].RunTime.UnixMilli()),
				Member: task.UnionTimerIDAndRunTime(tasks[i].Tid, tasks[i].RunTime.UnixMilli()),
			})
			// 默认 24小时过期
			pipe.Expire(ctx, tableName, time.Hour*24)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			log.Error(ctx, "pipe execute error.", "error", err)
			return err
		}
		return nil
	}); err != nil {
		log.Error(ctx, "save tasks to redis error.", "error", err)
		return err
	}
	log.Info(ctx, "save tasks to redis success.")
	return nil
}

func NewTaskCache(data *Data, conf *conf.Data) biz.TaskCache {
	return &taskCache{
		data:     data,
		confData: conf,
	}
}
