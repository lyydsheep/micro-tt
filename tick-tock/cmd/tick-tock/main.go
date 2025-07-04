package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"os"
	"tick-tock/configs"
	"tick-tock/pkg/log"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func newApp(gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	server, data := configs.NewConfig()

	app, cleanup, err := wireApp(server, data)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	log.Info(nil, "new app successfully",
		"service.id", id,
		"service.name", Name,
		"service.version", Version)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
