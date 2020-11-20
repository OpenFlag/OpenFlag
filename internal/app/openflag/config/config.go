package config

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/engine"

	"github.com/OpenFlag/OpenFlag/pkg/database"

	"github.com/OpenFlag/OpenFlag/pkg/monitoring/prometheus"
	"github.com/OpenFlag/OpenFlag/pkg/redis"

	"github.com/OpenFlag/OpenFlag/pkg/config"
	"github.com/OpenFlag/OpenFlag/pkg/log"
)

const (
	app       = "openflag"
	cfgFile   = "config.yaml"
	cfgPrefix = "openflag"
)

type (
	// Config represents application configuration struct.
	Config struct {
		Logger     Logger     `mapstructure:"logger"`
		Server     Server     `mapstructure:"server"`
		Database   Database   `mapstructure:"database"`
		Redis      Redis      `mapstructure:"redis"`
		Evaluation Evaluation `mapstructure:"evaluation"`
		Monitoring Monitoring `mapstructure:"monitoring"`
	}

	// Logger represents logger configuration struct.
	Logger struct {
		AccessLogger log.AccessLogger    `mapstructure:"access"`
		AppLogger    log.AppLogger       `mapstructure:"app"`
		Evaluation   engine.LoggerConfig `mapstructure:"evaluation"`
	}

	// Server represents server configuration struct.
	Server struct {
		Address         string        `mapstructure:"address"`
		RPCAddress      string        `mapstructure:"grpc-address"`
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

	// Evaluation represents evaluation configuration struct.
	Evaluation struct {
		EntityContextCacheExpiration time.Duration `mapstructure:"entity-context-cache-expiration"`
		UpdateFlagsCronPattern       string        `mapstructure:"update-flags-cron-pattern"`
	}

	// Monitoring represents monitoring configuration struct.
	Monitoring struct {
		Prometheus prometheus.Config `mapstructure:"prometheus"`
	}
)

// Validate validates Database struct.
func (d Database) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(
			&d.Driver,
			validation.In("postgres"),
		),
	)
}

// Validate validates Config struct.
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Database,
		),
	)
}

// Init initializes application configuration.
func Init() Config {
	var cfg Config

	config.Init(app, cfgFile, &cfg, defaultConfig, cfgPrefix)

	if err := cfg.Validate(); err != nil {
		logrus.Fatalf("failed to validate configurations: %s", err.Error())
	}

	return cfg
}
