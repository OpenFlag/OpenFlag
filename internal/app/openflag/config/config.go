package config

import (
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/evaluation"

	"github.com/OpenFlag/OpenFlag/pkg/database"

	"github.com/OpenFlag/OpenFlag/pkg/monitoring/prometheus"
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
		Database   Database   `mapstructure:"database" validate:"required"`
		Redis      Redis      `mapstructure:"redis" validate:"required"`
		Monitoring Monitoring `mapstructure:"monitoring" validate:"required"`
	}

	// Logger represents logger configuration struct.
	Logger struct {
		AccessLogger log.AccessLogger        `mapstructure:"access" validate:"required"`
		AppLogger    log.AppLogger           `mapstructure:"app" validate:"required"`
		Evaluation   evaluation.LoggerConfig `mapstructure:"evaluation" validate:"required"`
	}

	// Server represents server configuration struct.
	Server struct {
		Address         string        `mapstructure:"address" validate:"required"`
		ReadTimeout     time.Duration `mapstructure:"read-timeout" validate:"required"`
		WriteTimeout    time.Duration `mapstructure:"write-timeout" validate:"required"`
		GracefulTimeout time.Duration `mapstructure:"graceful-timeout" validate:"required"`
	}

	// Database represents database configuration struct.
	Database struct {
		Driver        string           `mapstructure:"driver" validate:"required"`
		MasterConnStr string           `mapstructure:"master-conn-string" validate:"required"`
		SlaveConnStr  string           `mapstructure:"slave-conn-string" validate:"required"`
		Options       database.Options `mapstructure:"options" validate:"required"`
	}

	// Redis represents Redis configuration struct.
	Redis struct {
		MasterAddress string        `mapstructure:"master-address" validate:"required"`
		SlaveAddress  string        `mapstructure:"slave-address" validate:"required"`
		Options       redis.Options `mapstructure:"options"`
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
