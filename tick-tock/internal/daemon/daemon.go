package daemon

import (
	"github.com/google/wire"
	"tick-tock/internal/biz"
	"tick-tock/internal/conf"
	"tick-tock/pkg/log"
)

var ProviderSet = wire.NewSet(NewHandles, NewTasks, NewServer)

// 需要注意配置顺序和任务顺序是否一致
func NewHandles(migrator *biz.MigratorUseCase, scheduler *biz.SchedulerUsecase) []Job {
	return []Job{
		migrator.Start,
		scheduler.Schedule,
	}
}

func NewTasks(conf *conf.Server, handles []Job) []Task {
	taskConf := conf.GetTask()
	configs := taskConf.GetTasks()
	if len(configs) != len(handles) {
		log.Fatal(nil, "task configs and handles length not equal.", "task configs length", len(configs), "handles length", len(handles))
	}

	tasks := make([]Task, 0, len(configs))
	for i := range configs {
		tasks = append(tasks, Task{
			Name:     configs[i].Name,
			Type:     taskType(configs[i].Type),
			Schedule: configs[i].Schedule,
			Handle:   handles[i],
		})
	}
	log.Info(nil, "new tasks successfully", "length of tasks", len(tasks))
	return tasks
}
