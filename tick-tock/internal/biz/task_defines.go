package biz

import (
	"context"
	"tick-tock/internal/constant"
	"time"
)

const TableNameTaskDefine = "task_defines"

// TaskDefine 任务定义表
type TaskDefine struct {
	ID              int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	Tid             string    `gorm:"column:tid;type:varchar(36);not null;comment:任务唯一标识（UUID）" json:"tid"`                                                      // 任务唯一标识（UUID）
	App             string    `gorm:"column:app;type:varchar(64);not null;comment:所属应用标识" json:"app"`                                                            // 所属应用标识
	Name            string    `gorm:"column:name;type:varchar(32);not null;comment:任务名称" json:"name"`                                                            // 任务名称
	Status          int32     `gorm:"column:status;type:tinyint;not null;default:1;comment:状态：1-active, 2-inactive" json:"status"`                               // 状态：1-active, 2-inactive
	Cron            string    `gorm:"column:cron;type:varchar(64);not null;comment:Cron表达式（标准5/6字段）" json:"cron"`                                                // Cron表达式（标准5/6字段）
	NotifyHTTPParam string    `gorm:"column:notify_http_param;type:json;comment:回调参数（JSON格式，如{"url":"","method":"POST","headers":{}}）" json:"notify_http_param"` // 回调参数（JSON格式，如{"url":"","method":"POST","headers":{}}）
	LastMigrateTime time.Time `gorm:"column:last_migrate_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"last_migrate_time"`
	CreateTime      time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
}

type NotifyHTTPParam struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Params  map[string]string `json:"params"`
	Body    string            `json:"body"`
}

// TableName TaskDefine's table name
func (*TaskDefine) TableName() string {
	return TableNameTaskDefine
}

type TaskDefineRepo interface {
	Create(ctx context.Context, taskDefine *TaskDefine) (*TaskDefine, error)
	Update(ctx context.Context, taskDefine *TaskDefine) (*TaskDefine, error)
	GetTaskDefineByTID(ctx context.Context, tID string) (*TaskDefine, error)
	GetTaskDefineByStatus(ctx context.Context, status constant.TaskDefineStatus) ([]*TaskDefine, error)
}
