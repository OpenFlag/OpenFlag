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
	app       = "openflag"
	cfgPath   = "config.yaml"
	cfgPrefix = "openflag"
)

type (
	// Config represents application configuration struct.
	Config struct {
		Logger     Logger     `mapstructure:"logger"`
		Server     Server     `mapstructure:"server"`
		Database   Database   `mapstructure:"database"`
		Redis      Redis      `mapstructure:"redis"`
		Monitoring Monitoring `mapstructure:"monitoring"`
	}

	// Logger represents logger configuration struct.
	Logger struct {
		AccessLogger log.AccessLogger        `mapstructure:"access"`
		AppLogger    log.AppLogger           `mapstructure:"app"`
		Evaluation   evaluation.LoggerConfig `mapstructure:"evaluation"`
	}

	// Server represents server configuration struct.
	Server struct {
		Address         string        `mapstructure:"address"`
		ReadTimeout     time.Duration `mapstructure:"read-timeout"`
		WriteTimeout    time.Duration `mapstructure:"write-timeout"`
		GracefulTimeout time.Duration `mapstructure:"graceful-timeout"`
	}

	// Database represents database configuration struct.
	Database struct {
		Driver        string           `mapstructure:"driver"`
		MasterConnStr string           `mapstructure:"master-conn-string"`
		SlaveConnStr  string           `mapstructure:"slave-conn-string"`
		Options       database.Options `mapstructure:"options"`
	}

	// Redis represents Redis configuration struct.
	Redis struct {
		MasterAddress string        `mapstructure:"master-address"`
		SlaveAddress  string        `mapstructure:"slave-address"`
		Options       redis.Options `mapstructure:"options"`
	}

	// Monitoring represents monitoring configuration struct.
	Monitoring struct {
		Prometheus prometheus.Config `mapstructure:"prometheus"`
	}
)

// Init initializes application configuration.
func Init() Config {
	var cfg Config

	config.Init(app, cfgPath, &cfg, defaultConfig, cfgPrefix)

	return cfg
}
