package data

import (
	"context"
	"tick-tock/internal/biz"
)

type taskDefineRepo struct {
	data *Data
}

func (repo *taskDefineRepo) Create(ctx context.Context, taskDefine *biz.TaskDefine) (*biz.TaskDefine, error) {
	q := repo.data.query.TaskDefine
	err := q.WithContext(ctx).Create(taskDefine)
	return taskDefine, err
}

func (repo *taskDefineRepo) Update(ctx context.Context, taskDefine *biz.TaskDefine) (*biz.TaskDefine, error) {
	q := repo.data.query.TaskDefine
	err := q.WithContext(ctx).Save(taskDefine)
	return taskDefine, err
}

func (repo *taskDefineRepo) GetTaskDefineByID(ctx context.Context, id int64) (*biz.TaskDefine, error) {
	q := repo.data.query.TaskDefine
	return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

func NewTaskDefineRepo(data *Data) biz.TaskDefineRepo {
	return &taskDefineRepo{
		data: data,
	}
}
