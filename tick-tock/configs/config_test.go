package configs

import (
	"flag"
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	flag.Parse()
	server, data := NewConfig()
	fmt.Println(data.Scheduler.LockDuration.AsDuration())
	_, _ = server, data
}
