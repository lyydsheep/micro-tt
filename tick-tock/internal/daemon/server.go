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

func NewServer(serverConf *conf.Server, tasks []Task) *Server {
	return &Server{
		serverConf: serverConf,
		tasks:      tasks,
	}
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
