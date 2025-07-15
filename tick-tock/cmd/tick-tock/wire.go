//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"tick-tock/internal/biz"
	"tick-tock/internal/conf"
	"tick-tock/internal/daemon"
	"tick-tock/internal/data"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, daemon.ProviderSet, newApp))
}
