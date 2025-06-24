package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"time"
)

func main() {
	httpSrv := http.NewServer(http.Address(":8080"))
	app := kratos.New(kratos.Name("tick-tock"), kratos.Server(httpSrv))

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				log.Info("tick")
			}
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
