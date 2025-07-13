package biz

import (
	"context"
	"time"
)

const TableNameTask = "task"

// Task mapped from table <task>
type Task struct {
	ID         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	App        string    `gorm:"column:app;type:varchar(64);not null" json:"app"`
	Tid        string    `gorm:"column:tid;type:varchar(36);not null" json:"tid"`
	Output     string    `gorm:"column:output;type:text;not null;comment:执行结果" json:"output"`                              // 执行结果
	RunTime    time.Time `gorm:"column:run_time;type:timestamp;not null;comment:执行时间" json:"run_time"`                     // 执行时间
	CostTime   int64     `gorm:"column:cost_time;type:bigint;not null;comment:执行耗时（毫秒）" json:"cost_time"`                  // 执行耗时（毫秒）
	Status     int32     `gorm:"column:status;type:tinyint;not null;comment:当前状态 未执行 0 执行中 1 执行成功 2 执行失败 3" json:"status"` // 当前状态 未执行 0 执行中 1 执行成功 2 执行失败 3
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName Task's table name
func (*Task) TableName() string {
	return TableNameTask
}

type TaskRepo interface {
	Create(ctx context.Context, task *Task) (*Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	GetTaskByRunTime(ctx context.Context, startTime time.Time, endTime time.Time, status int32) ([]*Task, error)
	UpdateByTIDAndRuntime(ctx context.Context, tid string, runTime time.Time, updates map[string]any) error
}

type TaskCache interface {
	GetValByStartAndEnd(ctx context.Context, key string, startUnixMilli int64, endUnixMilli int64) ([]*Task, error)
	SaveTasks(ctx context.Context, tasks []*Task) error
}
