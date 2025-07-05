package data

import (
	"context"
	"tick-tock/internal/data/gen"
	"tick-tock/pkg/log"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTaskRepo, NewTaskDefineRepo)

// 传递 transaction db
const contextVal = "transaction"

// Data .
type Data struct {
	query *gen.Query
}

// NewData .
func NewData(query *gen.Query) (*Data, func(), error) {
	cleanup := func() {
		log.Info(nil, "closing the data resources.")
	}
	return &Data{}, cleanup, nil
}

func (d *Data) Txn(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.query.Transaction(func(tx *gen.Query) error {
		ctx = context.WithValue(ctx, contextVal, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gen.Query {
	q, ok := ctx.Value(contextVal).(*gen.Query)
	if ok {
		return q
	}
	return d.query
}
