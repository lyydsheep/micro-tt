package biz

import (
	"context"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
// 注入 useCase
var ProviderSet = wire.NewSet(NewMigrator)

type Transaction interface {
	Txn(ctx context.Context, fn func(ctx context.Context) error) error
}
