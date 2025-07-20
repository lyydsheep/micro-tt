package data

import (
	"context"
	"tick-tock/internal/biz"
	"tick-tock/internal/constant"
	"time"
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
	_, err := q.WithContext(ctx).Where(q.ID.Eq(taskDefine.ID)).Updates(map[string]any{
		"last_migrate_time": taskDefine.LastMigrateTime,
		"update_time":       time.Now().UTC(),
	})
	return taskDefine, err
}

func (repo *taskDefineRepo) GetTaskDefineByTID(ctx context.Context, tID string) (*biz.TaskDefine, error) {
	q := repo.data.query.TaskDefine
	return q.WithContext(ctx).Where(q.Tid.Eq(tID)).First()
}

func NewTaskDefineRepo(data *Data) biz.TaskDefineRepo {
	return &taskDefineRepo{
		data: data,
	}
}
