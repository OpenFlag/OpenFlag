package config

import (
	"github.com/OpenFlag/OpenFlag/pkg/config"
)

const (
	cfgPath   = "config.yaml"
	cfgPrefix = "openflag"
)

type (
	Config struct {
		Logger Logger `mapstructure:"logger" validate:"required"`
	}

	Logger struct {
		Level string `mapstructure:"level" validate:"required"`
	}
)

func Init() Config {
	var cfg Config

	config.Init(cfgPath, &cfg, defaultConfig, cfgPrefix)

	return cfg
}
