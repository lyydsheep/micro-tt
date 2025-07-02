package configs

import (
	"flag"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"tick-tock/internal/conf"
	"tick-tock/pkg/log"
)

var (
	configPath = flag.String("config", "../../configs/config.yaml", "config file path. eg: ./config.yaml")
)

func NewConfig() (*conf.Server, *conf.Data) {
	cfg := config.New(config.WithSource(file.NewSource(*configPath)))
	defer cfg.Close()
	if err := cfg.Load(); err != nil {
		log.Fatal(nil, "load config err.", "error", err)
	}

	var c conf.Bootstrap
	if err := cfg.Scan(&c); err != nil {
		log.Fatal(nil, "scan config err.", "error", err)
	}
	log.Debug(nil, "new config successfully.", "config", &c)
	return c.Server, c.Data
}
