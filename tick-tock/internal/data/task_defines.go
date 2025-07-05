package data

import (
	"context"
	"tick-tock/internal/biz"
	"tick-tock/internal/constant"
)

type taskDefineRepo struct {
	data *Data
}

func (repo *taskDefineRepo) GetTaskDefineByStatus(ctx context.Context, status constant.TaskDefineStatus) ([]*biz.TaskDefine, error) {
	q := repo.data.query.TaskDefine
	return q.WithContext(ctx).Where(q.Status.Eq(status.ToInt32())).Find()
}

func (repo *taskDefineRepo) Create(ctx context.Context, taskDefine *biz.TaskDefine) (*biz.TaskDefine, error) {
	q := repo.data.DB(ctx).TaskDefine
	err := q.WithContext(ctx).Create(taskDefine)
	return taskDefine, err
}

func (repo *taskDefineRepo) Update(ctx context.Context, taskDefine *biz.TaskDefine) (*biz.TaskDefine, error) {
	q := repo.data.DB(ctx).TaskDefine
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
