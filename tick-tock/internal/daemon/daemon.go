package daemon

import (
	"context"
	"tick-tock/pkg/log"
	"time"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewHandles, NewServer)

func NewHandles() []Job {
	return []Job{
		// mock migrator
		func(ctx context.Context) {
			log.Info(nil, "migrator working.")
		},
		// mock scheduler
		func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					log.Info(ctx, "scheduler stop.")
				default:

				}
				log.Info(nil, "scheduler working.")
				time.Sleep(time.Second * 10)
			}
		},
	}
}
