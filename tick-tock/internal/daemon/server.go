package daemon

import (
	"context"
	"tick-tock/internal/conf"
	"tick-tock/pkg/log"
)

type Server struct {
	cancel     context.CancelFunc
	serverConf *conf.Server
	tasks      []Task
}

func NewServer(serverConf *conf.Server, jobs []Job) *Server {
	server := &Server{
		serverConf: serverConf,
	}
	server.tasks = server.newTasks(jobs)
	return server
}

func (server *Server) newTasks(handles []Job) []Task {
	taskConf := server.serverConf.GetTask()
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
	return tasks
}

func (server *Server) Start(ctx context.Context) error {
	log.Info(ctx, "daemon start.")
	ctx, cancel := context.WithCancel(ctx)
	server.cancel = cancel
	for _, task := range server.tasks {
		go task.Run(ctx)
	}
	return nil
}

func (server *Server) Stop(ctx context.Context) error {
	log.Info(ctx, "daemon stop.")
	server.cancel()
	return nil
}
