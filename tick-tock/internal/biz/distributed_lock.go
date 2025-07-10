package biz

import (
	"context"
	"time"
)

type Lock interface {
	Lock(ctx context.Context, key string, expire time.Duration) (bool, string, error)
	RenewLock(ctx context.Context, key string, val string, newTTL time.Duration) error
}
