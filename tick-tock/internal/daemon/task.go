package daemon

import (
	"context"
	"github.com/gorhill/cronexpr"
	"tick-tock/pkg/log"
	"time"
)

type taskType string

const (
	Cron taskType = "cron"
	Once taskType = "once"
)

type Job func(ctx context.Context)

type Task struct {
	Name     string
	Type     taskType
	Schedule string
	Handle   Job
}

func (t *Task) Run(ctx context.Context) {
	switch {
	case t.Type == Cron:
		t.cron(ctx)
	case t.Type == Once:
		t.once(ctx)
	default:
		log.Warn(nil, "task type not support", "type", t.Type)
	}
}

func (t *Task) once(ctx context.Context) {
	log.Debug(ctx, "task start", "name", t.Name, "type", t.Type, "schedule", t.Schedule)
	go func() {
		// 定时执行
		waitTime, err := time.ParseDuration(t.Schedule)
		if err != nil {
			log.Error(ctx, "parse duration error.", "error", err, "taskName", t.Name)
			return
		}
		if err = t.run(ctx, waitTime); err != nil {
			log.Error(ctx, "task run error.", "error", err, "taskName", t.Name)
		}
	}()
}

func (t *Task) cron(ctx context.Context) {
	// 周期性执行
	log.Info(ctx, "task run", "taskName", t.Name, "taskType", t.Type, "taskSchedule", t.Schedule)
	go func() {
		expr, err := cronexpr.Parse(t.Schedule)
		if err != nil {
			log.Fatal(ctx, "parse cron expression error.", "error", err, "taskName", t.Name)
		}
		// 立即执行一次，然后等待下一个周期
		t.Handle(ctx)
		for {
			wait := expr.Next(time.Now()).Sub(time.Now())
			if err = t.run(ctx, wait); err != nil {
				log.Error(ctx, "task run error.", "error", err, "taskName", t.Name)
				return
			}
		}
	}()
}

func (t *Task) run(ctx context.Context, wait time.Duration) error {
	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-timer.C:
		log.Info(ctx, "task run", "taskName", t.Name)
		t.Handle(ctx)
	case <-ctx.Done():
		log.Info(ctx, "task stop", "taskName", t.Name)
		return ctx.Err()
	}
	return nil
}
