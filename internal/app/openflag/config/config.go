package config

import (
	"time"

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

	Monitoring struct {
		Prometheus Prometheus `mapstructure:"prometheus" validate:"required"`
	}

	Prometheus struct {
		Enabled bool   `mapstructure:"enabled"`
		Address string `mapstructure:"address" validate:"required"`
	}
)

func Init() Config {
	var cfg Config

	config.Init(cfgPath, &cfg, defaultConfig, cfgPrefix)

	return cfg
}
