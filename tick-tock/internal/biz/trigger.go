package biz

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/sync/errgroup"
	"tick-tock/internal/conf"
	"tick-tock/internal/constant"
	"tick-tock/pkg/log"
	"tick-tock/util/task"
	"time"
)

type TriggerUsecase struct {
	conf     *conf.Data
	pool     ants.Pool
	cache    TaskCache
	poll     *ants.Pool
	taskRepo TaskRepo
}

func (uc *TriggerUsecase) Work(ctx context.Context, tableName string, ack func()) error {
	log.Info(ctx, "start trigger work.", "tableName", tableName)
	// table 中存储的是一分钟内的任务
	// TODO
	ticker := time.NewTicker(uc.conf.Trigger.PollInterval.AsDuration())
	defer ticker.Stop()

	var eg errgroup.Group
	startTime, err := task.GetRuntimeMinute(ctx, tableName)
	if err != nil {
		log.Error(ctx, "get runtime minute error.", "err", err, "tableName", tableName)
		return err
	}
	endTime := startTime.Add(time.Minute)
	rangeGap := uc.conf.Trigger.RangeGap.AsDuration()
	if rangeGap <= 0 {
		log.Warn(ctx, "range gap is not greater zero, use default value.", "rangeGap", rangeGap)
		rangeGap = time.Second
	}
	for range ticker.C {
		select {
		case <-ctx.Done():
			log.Info(ctx, "context cancel, trigger stop.", "tableName", tableName)
			return nil
		default:
			eg.Go(func() error {
				return uc.handleSlice(ctx, tableName, startTime.UnixMilli(),
					startTime.Add(rangeGap).UnixMilli())
			})
			if startTime.Add(rangeGap); !startTime.Before(endTime) {
				log.Info(ctx, "trigger had gone all slices.", "tableName", tableName)
				break
			}
		}
	}
	if err = eg.Wait(); err != nil {
		log.Error(ctx, "trigger error.", "error", err, "tableName", tableName)
		return err
	}
	return nil
}

func (uc *TriggerUsecase) handleSlice(ctx context.Context, tableName string, startUnixMilli int64, endUnixMilli int64) error {
	start, end := time.UnixMilli(startUnixMilli), time.UnixMilli(endUnixMilli)
	log.Info(ctx, "handle slice.", "tableName", tableName, "start", start, "end", end)
	// 先走 cache
	tasks, err := uc.cache.GetValByStartAndEnd(ctx, tableName, startUnixMilli, endUnixMilli)
	if err != nil {
		log.Error(ctx, "fail to get tasks from cache.", "error", err, "tableName", tableName, "start", start, "end", end)
	}
	// 再走 db
	// 根据 runtime 和 status 过滤任务
	tasks, err = uc.taskRepo.GetTaskByRunTime(ctx, start, end)
	if err != nil {
		log.Error(ctx, "fail to get tasks from db.", "error", err, "tableName", tableName, "start", start, "end", end)
		return err
	}

	// 在根据 taskid 过滤出属于这个 bucket 的任务
	bucketId, err := task.GetBucketID(ctx, tableName)
	for i := range tasks {
		if tasks[i].ID%int64(uc.conf.Scheduler.BucketCount) == int64(bucketId) && tasks[i].Status == constant.TaskInit.ToInt32() {
			// pool 过载或者已关闭则会报错
			if err = uc.pool.Submit(func() {

			}); err != nil {

			}
		}
	}

	return nil
}
