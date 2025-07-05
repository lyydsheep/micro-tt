package data

import (
	"context"
	"tick-tock/internal/biz"
	"time"
)

type taskRepo struct {
	data *Data
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

func (repo *taskRepo) GetTaskByIDAndRunTime(ctx context.Context, id int64, runTime time.Time) (*biz.Task, error) {
	q := repo.data.query.Task
	task, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Where(q.RunTime.Eq(runTime)).First()
	return task, err
}

func NewTaskRepo(data *Data) biz.TaskRepo {
	return &taskRepo{
		data: data,
	}
}
