package config

import (
	"bytes"
	"log"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

func Init(path string, cfg interface{}, defaultConfig string, prefix string) interface{} {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader([]byte(defaultConfig))); err != nil {
		log.Fatalf("error loading default configs: %s", err.Error())
	}

	v.SetConfigFile(path)
	v.SetEnvPrefix(prefix)
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	err := v.MergeInConfig()
	if err != nil {
		log.Println("no config file found. Using defaults and environment variables")
	}

	if err := v.UnmarshalExact(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config into struct: %s", err.Error())
	}

	if err := validate(cfg); err != nil {
		log.Fatalf("failed to validate configuration: %s", err.Error())
	}

	return cfg
}

func validate(cfg interface{}) error {
	return validator.New().Struct(cfg)
}
