package data

import (
	"context"
	"tick-tock/internal/biz"
	"time"
)

type taskRepo struct {
	data *Data
}

func (repo *taskRepo) UpdateByTIDAndRuntime(ctx context.Context, tid string, runTime time.Time, updates map[string]any) error {
	q := repo.data.DB(ctx).Task
	_, err := q.WithContext(ctx).Where(q.Tid.Eq(tid), q.RunTime.Eq(runTime)).Updates(updates)
	return err
}

func (repo *taskRepo) Create(ctx context.Context, task *biz.Task) (*biz.Task, error) {
	q := repo.data.DB(ctx).Task
	err := q.WithContext(ctx).Create(task)
	return task, err
}

func (repo *taskRepo) Update(ctx context.Context, task *biz.Task) (*biz.Task, error) {
	q := repo.data.DB(ctx).Task
	err := q.WithContext(ctx).Save(task)
	return task, err
}

func (repo *taskRepo) GetTaskByRunTime(ctx context.Context, startTime time.Time, endTime time.Time, status int32) ([]*biz.Task, error) {
	q := repo.data.query.Task
	tasks, err := q.WithContext(ctx).Where(q.RunTime.Gte(startTime),
		q.RunTime.Lt(endTime), q.Status.Eq(status)).Find()
	return tasks, err
}

func NewTaskRepo(data *Data) biz.TaskRepo {
	return &taskRepo{
		data: data,
	}
}
