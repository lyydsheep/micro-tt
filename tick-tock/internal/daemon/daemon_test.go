package daemon

import (
	"context"
	"testing"
	"tick-tock/configs"
	"time"
)

func TestNewServer(t *testing.T) {
	server, _ := configs.NewConfig()
	daemon := NewServer(server, NewHandles())
	ctx := context.Background()
	daemon.Start(ctx)
	time.Sleep(time.Minute * 2)
	daemon.Stop(ctx)
}
