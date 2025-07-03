package data

import (
	"tick-tock/internal/data/gen"
	"tick-tock/pkg/log"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

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
