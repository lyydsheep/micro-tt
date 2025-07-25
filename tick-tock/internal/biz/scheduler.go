package biz

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"tick-tock/internal/conf"
	"tick-tock/pkg/log"
	"tick-tock/util/task"
	"time"
)

type SchedulerUsecase struct {
	lock    Lock
	conf    *conf.Data
	pool    *ants.Pool
	trigger *TriggerUsecase
}

func NewSchedulerUsecase(lock Lock, conf *conf.Data, trigger *TriggerUsecase) *SchedulerUsecase {
	pool, err := ants.NewPool(int(conf.Scheduler.WorkerPoolSize), ants.WithNonblocking(true))
	if err != nil {
		log.Fatal(nil, "new ants pool failed.", "error", err)
	}
	return &SchedulerUsecase{
		lock:    lock,
		conf:    conf,
		pool:    pool,
		trigger: trigger,
	}
}

func (uc *SchedulerUsecase) Schedule(ctx context.Context) {
	// 一分钟内有多个 bucket
	// 根据配置对 Redis 中当前分钟内的任务进行扫描  左闭右开
	// eg: [2022-01-01 00:00:00, 2022-01-01 00:01:00)
	log.Info(ctx, "start schedule")
	ticker := time.NewTicker(uc.conf.Scheduler.PollInterval.AsDuration())
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-ctx.Done():
			log.Info(ctx, "context cancel, scheduler stop.")
			return
		default:
			now := time.Now().UTC()
			for i := range uc.conf.Scheduler.BucketCount {
				// 重试上一分钟的任务
				if err := uc.pool.Submit(func() {
					uc.handleBucket(ctx, now.Add(-time.Minute), i)
				}); err != nil {
					log.Error(ctx, "submit task error.", "error", err, "bucketID", i, "startTime", now.Add(-time.Minute))
				}
				// 处理当前分钟的任务
				if err := uc.pool.Submit(func() {
					uc.handleBucket(ctx, now, i)
				}); err != nil {
					log.Error(ctx, "submit task error.", "error", err, "bucketID", i, "startTime", now)
				}
			}
		}
	}
}

func (uc *SchedulerUsecase) handleBucket(ctx context.Context, now time.Time, bucketID int32) {
	tableName := task.GetTableName(ctx, uc.conf, now, int64(bucketID))
	// 抢锁
	key := task.GetLockKey(ctx, uc.conf, tableName)
	res, val, err := uc.lock.Lock(ctx, key, uc.conf.Scheduler.LockDuration.AsDuration())
	if err != nil {
		log.Error(ctx, "lock failed.", "error", err, "key", key)
		return
	}
	if !res {
		log.Info(ctx, "distributed lock has been taken.", "key", key)
		return
	}
	log.Info(ctx, "get distributed lock success.", "key", key, "val", val)
	ack := func() {
		// 任务成功完成，调用此函数
		// 对分布式锁进行续期
		if err := uc.lock.RenewLock(ctx, key, val, uc.conf.Scheduler.RenewLockDuration.AsDuration()); err != nil {
			log.Error(ctx, "renew lock failed.", "error", err, "key", key)
			return
		}
		log.Info(ctx, "renew lock success.", "key", key)
	}
	if err = uc.trigger.Work(ctx, tableName, ack); err != nil {
		log.Error(ctx, "trigger work failed.", "error", err, "tableName", tableName)
	}
}
