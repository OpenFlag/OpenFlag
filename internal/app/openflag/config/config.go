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
	// Config represents application configuration struct.
	Config struct {
		Logger     Logger     `mapstructure:"logger" validate:"required"`
		Server     Server     `mapstructure:"server" validate:"required"`
		Postgres   Postgres   `mapstructure:"postgres" validate:"required"`
		Redis      Redis      `mapstructure:"redis" validate:"required"`
		Monitoring Monitoring `mapstructure:"monitoring" validate:"required"`
	}

	// Logger represents logger configuration struct.
	Logger struct {
		AccessLogger log.AccessLogger `mapstructure:"access" validate:"required"`
		AppLogger    log.AppLogger    `mapstructure:"app" validate:"required"`
	}

	// Server represents server configuration struct.
	Server struct {
		Address         string        `mapstructure:"address" validate:"required"`
		ReadTimeout     time.Duration `mapstructure:"read-timeout" validate:"required"`
		WriteTimeout    time.Duration `mapstructure:"write-timeout" validate:"required"`
		GracefulTimeout time.Duration `mapstructure:"graceful-timeout" validate:"required"`
	}

	// Postgres represents PostgreSQL configuration struct.
	Postgres struct {
		Master postgres.Config `mapstructure:"master" validate:"required"`
		Slave  postgres.Config `mapstructure:"slave" validate:"required"`
	}

	// Redis represents Redis configuration struct.
	Redis struct {
		Master redis.Config `mapstructure:"master" validate:"required"`
		Slave  redis.Config `mapstructure:"slave" validate:"required"`
	}

	// Monitoring represents monitoring configuration struct.
	Monitoring struct {
		Prometheus prometheus.Config `mapstructure:"prometheus" validate:"required"`
	}
)

// Init initializes application configuration.
func Init() Config {
	var cfg Config

	config.Init(cfgPath, &cfg, defaultConfig, cfgPrefix)

	return cfg
}
