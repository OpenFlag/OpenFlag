package config

import (
	"bytes"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

// Init initializes a config struct using default, file, and environment variables.
func Init(path string, cfg interface{}, defaultConfig string, prefix string) interface{} {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader([]byte(defaultConfig))); err != nil {
		logrus.Fatalf("error loading default configs: %s", err.Error())
	}

	v.SetConfigFile(path)
	v.SetEnvPrefix(prefix)
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	err := v.MergeInConfig()
	//nolint:staticcheck
	if err != nil {
		//logrus.Warn("no config file found. Using defaults and environment variables")
	}

	if err := v.UnmarshalExact(&cfg); err != nil {
		logrus.Fatalf("failed to unmarshal config into struct: %s", err.Error())
	}

	if err := validate(cfg); err != nil {
		logrus.Fatalf("failed to validate configuration: %s", err.Error())
	}

	return cfg
}

func validate(cfg interface{}) error {
	return validator.New().Struct(cfg)
}
