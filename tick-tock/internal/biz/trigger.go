package biz

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"tick-tock/internal/conf"
)

type TriggerUsecase struct {
	conf *conf.Data
	pool ants.Pool
}

func (uc *TriggerUsecase) Work(ctx context.Context, tableName string, ack func()) error {
	// table 中存储的是一分钟内的任务
	// TODO

	return nil
}
