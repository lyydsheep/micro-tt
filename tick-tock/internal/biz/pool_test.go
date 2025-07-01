package biz

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"testing"
	"tick-tock/pkg/log"
	"time"
)

func TestAnts(t *testing.T) {
	pool, _ := ants.NewPool(100, ants.WithExpiryDuration(time.Second),
		ants.WithNonblocking(true), ants.WithPreAlloc(true))
	if err := pool.Submit(func() {
		time.Sleep(time.Second)
		fmt.Println(time.Now().UTC())
		fmt.Println("hello world")
	}); err != nil {
		log.Fatal(nil, "submit task error: %v", err)
	}
	time.Sleep(time.Second * 2)
}
