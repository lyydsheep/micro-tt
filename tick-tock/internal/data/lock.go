package data

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"tick-tock/internal/biz"
	"time"
)

type lock struct {
	data *Data
}

func (l *lock) Lock(ctx context.Context, key string, expire time.Duration) (bool, string, error) {
	val := uuid.New().String()
	res, err := l.data.redis.SetNX(ctx, key, val, expire).Result()
	return res, val, err
}

func (l *lock) RenewLock(ctx context.Context, key string, val string, newTTL time.Duration) error {
	luaScript := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
	`
	res, err := l.data.redis.Eval(ctx, luaScript, []string{key}, val, newTTL.Milliseconds()).Result()
	if err != nil {
		return err
	}
	if res.(int64) == 0 {
		return errors.New("lock is not owned")
	}
	return nil
}

func NewLock(data *Data) biz.Lock {
	return &lock{
		data: data,
	}
}
