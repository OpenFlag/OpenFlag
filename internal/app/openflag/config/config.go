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
		Host               string        `mapstructure:"host" validate:"required"`
		Port               int           `mapstructure:"port" validate:"required"`
		Username           string        `mapstructure:"user" validate:"required"`
		Password           string        `mapstructure:"pass" validate:"required"`
		DBName             string        `mapstructure:"dbname" validate:"required"`
		ConnectTimeout     time.Duration `mapstructure:"connect-timeout" validate:"required"`
		ConnectionLifetime time.Duration `mapstructure:"connection-lifetime" validate:"required"`
		MaxOpenConnections int           `mapstructure:"max-open-connections"`
		MaxIdleConnections int           `mapstructure:"max-idle-connections"`
	}

	Redis struct {
		Address         string        `mapstructure:"address" validate:"required"`
		PoolSize        int           `mapstructure:"pool-size"`
		MinIdleConns    int           `mapstructure:"min-idle-conns"`
		DialTimeout     time.Duration `mapstructure:"dial-timeout"`
		ReadTimeout     time.Duration `mapstructure:"read-timeout"`
		WriteTimeout    time.Duration `mapstructure:"write-timeout"`
		PoolTimeout     time.Duration `mapstructure:"pool-timeout"`
		IdleTimeout     time.Duration `mapstructure:"idle-timeout"`
		MaxRetries      int           `mapstructure:"max-retries"`
		MinRetryBackoff time.Duration `mapstructure:"min-retry-backoff"`
		MaxRetryBackoff time.Duration `mapstructure:"max-retry-backoff"`
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
