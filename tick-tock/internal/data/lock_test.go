package data

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	data, _, err := NewData(nil, redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:36379",
	}))
	if err != nil {
		t.Error(err)
		return
	}
	l := NewLock(data)
	ctx := context.Background()
	res, val, err := l.Lock(ctx, "test", time.Second*10)
	if err != nil {
		t.Error(err)
		return
	}
	if !res {
		t.Log("lock failed")
		return
	}
	time.Sleep(time.Second * 8)
	go func() {
		res, _, err := l.Lock(ctx, "test", time.Second*10)
		if err != nil {
			t.Error(err)
			return
		}
		if !res {
			t.Log("lock failed")
			return
		}
	}()
	if err = l.RenewLock(ctx, "test", val, time.Hour); err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 10)
}
