package config

import (
	"time"

	"github.com/OpenFlag/OpenFlag/pkg/monitoring/prometheus"
	"github.com/OpenFlag/OpenFlag/pkg/postgres"
	"github.com/OpenFlag/OpenFlag/pkg/redis"

	"github.com/OpenFlag/OpenFlag/pkg/config"
	"github.com/OpenFlag/OpenFlag/pkg/log"
)

const (
	cfgPath   = "config.yaml"
	cfgPrefix = "openflag"
)

type (
	Config struct {
		Logger     Logger     `mapstructure:"logger" validate:"required"`
		Server     Server     `mapstructure:"server" validate:"required"`
		Postgres   Postgres   `mapstructure:"postgres" validate:"required"`
		Redis      Redis      `mapstructure:"redis" validate:"required"`
		Monitoring Monitoring `mapstructure:"monitoring" validate:"required"`
	}

	Logger struct {
		AccessLogger log.AccessLogger `mapstructure:"access" validate:"required"`
		AppLogger    log.AppLogger    `mapstructure:"app" validate:"required"`
	}

	Server struct {
		Address         string        `mapstructure:"address" validate:"required"`
		ReadTimeout     time.Duration `mapstructure:"read-timeout" validate:"required"`
		WriteTimeout    time.Duration `mapstructure:"write-timeout" validate:"required"`
		GracefulTimeout time.Duration `mapstructure:"graceful-timeout" validate:"required"`
	}

	Postgres struct {
		Master postgres.Config `mapstructure:"master" validate:"required"`
		Slave  postgres.Config `mapstructure:"slave" validate:"required"`
	}

	Redis struct {
		Master redis.Config `mapstructure:"master" validate:"required"`
		Slave  redis.Config `mapstructure:"slave" validate:"required"`
	}

	Monitoring struct {
		Prometheus prometheus.Config `mapstructure:"prometheus" validate:"required"`
	}
)

func Init() Config {
	var cfg Config

	config.Init(cfgPath, &cfg, defaultConfig, cfgPrefix)

	return cfg
}
