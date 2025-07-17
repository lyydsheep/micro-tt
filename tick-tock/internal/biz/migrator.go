package biz

import (
	"context"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
	"tick-tock/internal/conf"
	"tick-tock/internal/constant"
	"tick-tock/pkg/log"
	"time"
)

type MigratorUseCase struct {
	data           *conf.Data
	taskDefineRepo TaskDefineRepo
	taskRepo       TaskRepo
	txnManager     Transaction
	taskCache      TaskCache
}

func NewMigrator(data *conf.Data, taskDefineRepo TaskDefineRepo, taskRepo TaskRepo, tm Transaction) *MigratorUseCase {
	return &MigratorUseCase{
		data:           data,
		taskDefineRepo: taskDefineRepo,
		taskRepo:       taskRepo,
		txnManager:     tm,
	}
}

// Start 启动一次数据迁移
func (m *MigratorUseCase) Start(ctx context.Context) {
	log.Info(ctx, "migrator start.")
	// 获取所有 有效的 任务定义
	taskDefines, err := m.taskDefineRepo.GetTaskDefineByStatus(ctx, constant.TaskDefineActive)
	if err != nil {
		log.Error(ctx, "get task define error.", "error", err)
		return
	}
	log.Info(ctx, "get task define success.", "task_define_count", len(taskDefines))
	var eg errgroup.Group
	for _, taskDefine := range taskDefines {
		eg.Go(func() error {
			// 迁移一个步长的数据
			return m.migrator(ctx, taskDefine, m.data.Migrator.MigrateStep.AsDuration())
		})
	}
	if err = eg.Wait(); err != nil {
		log.Error(ctx, "migrator error.", "error", err)
	}
}

func (m *MigratorUseCase) migrator(ctx context.Context, taskDefine *TaskDefine, step time.Duration) error {
	// 必须是有效的任务
	if taskDefine.Status != constant.TaskDefineActive.ToInt32() {
		log.Warn(ctx, "task define status is not active.", "task_define_id", taskDefine.ID)
		return nil
	}
	log.Info(ctx, "migrator working.", "task_define_id", taskDefine.ID)

	// 生成一个步长内的定时任务
	start, now := taskDefine.LastMigrateTime, time.Now().UTC()
	// 如果开始时间早于当前时间，则从当前时间开始
	// 避免出现[start, now]之间的任务
	if start.Before(now) {
		start = now
		// 一次性生成两个步长时间的任务
		step <<= 1
		log.Info(ctx, "start time is before now, generate two steps in one time.", "task_define_id", taskDefine.ID)
	}
	end := start.Add(step)
	// 生成[start, end)时间范围内的任务
	tasks, err := m.generateTask(ctx, taskDefine, start, end)
	if err != nil {
		log.Error(ctx, "generate task error.", "error", err, "task_define_id", taskDefine.ID)
		return err
	}
	log.Info(ctx, "generate task success.", "task_define_id", taskDefine.ID, "start", start, "end", end)

	// 将定时任务插入数据库
	taskDefine.LastMigrateTime = end
	// 执行事务：1. 保存任务 2. 更新 任务定义的生成任务起始时间
	if err = m.txnManager.Txn(ctx, func(ctx context.Context) error {
		return m.saveTasksAndUpdateTaskDefine(ctx, taskDefine, tasks)
	}); err != nil {
		log.Error(ctx, "save tasks and update task define error.", "error", err, "task_define_id", taskDefine.ID)
		return err
	}

	// 将定时任务插入 Redis
	if err = m.saveTasks(ctx, tasks); err != nil {
		log.Error(ctx, "save tasks error.", "error", err, "task_define_id", taskDefine.ID)
		return err
	}

	return nil
}

func (m *MigratorUseCase) generateTask(ctx context.Context, taskDefine *TaskDefine, start, end time.Time) ([]*Task, error) {
	schedule, err := cron.ParseStandard(taskDefine.Cron)
	if err != nil {
		log.Error(ctx, "parse cron error.", "error", err, "cron", taskDefine.Cron, "task_define_id", taskDefine.ID)
		return nil, err
	}
	// 时间范围 [start, end)
	tasks := make([]*Task, 0)
	for next := schedule.Next(start); next.Before(end); next = schedule.Next(next) {
		tasks = append(tasks, &Task{
			App:        taskDefine.App,
			Tid:        taskDefine.Tid,
			Status:     constant.TaskInit.ToInt32(),
			RunTime:    next,
			CreateTime: time.Now().UTC(),
			UpdateTime: time.Now().UTC(),
		})
	}
	return tasks, nil
}

func (m *MigratorUseCase) saveTasksAndUpdateTaskDefine(ctx context.Context, taskDefine *TaskDefine, tasks []*Task) error {
	if _, err := m.taskDefineRepo.Update(ctx, taskDefine); err != nil {
		log.Error(ctx, "update task define error.", "error", err, "task_define_id", taskDefine.ID)
		return err
	}
	for _, task := range tasks {
		if _, err := m.taskRepo.Create(ctx, task); err != nil {
			log.Error(ctx, "create task error.", "error", err, "task_define_id", taskDefine.ID)
			return err
		}
	}
	return nil
}

func (m *MigratorUseCase) saveTasks(ctx context.Context, tasks []*Task) error {
	return m.taskCache.SaveTasks(ctx, tasks)
}
