package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Init initializes a config struct using default, file, and environment variables.
func Init(app string, file string, cfg interface{}, defaultConfig string, prefix string) interface{} {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader([]byte(defaultConfig))); err != nil {
		logrus.Fatalf("error loading default configs: %s", err.Error())
	}

	v.SetConfigFile(file)
	v.SetEnvPrefix(prefix)
	v.AddConfigPath(fmt.Sprintf("/etc/%s/", app))
	v.AddConfigPath(fmt.Sprintf("$HOME/.%s", app))
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	switch err := v.MergeInConfig(); err.(type) {
	case nil:
	case *os.PathError:
		logrus.Warn("no config file found. Using defaults and environment variables")
	default:
		logrus.Warnf("failed to load config file: %s", err.Error())
	}

	if err := v.UnmarshalExact(&cfg); err != nil {
		logrus.Fatalf("failed to unmarshal config into struct: %s", err.Error())
	}

	return cfg
}
