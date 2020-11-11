package log

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// AppLogger represents a struct for application logger configurations.
type AppLogger struct {
	Level      string `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`
	StdOut     bool   `mapstructure:"stdout"`
}

// SetupLogger setup application logger.
func SetupLogger(cfg AppLogger) {
	logLevel, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logLevel = logrus.DebugLevel
	}

	logrus.SetLevel(logLevel)

	if !cfg.StdOut {
		logrus.SetOutput(ioutil.Discard)

		logrus.AddHook(NewRotateFileHook(RotateFileConfig{
			Filename:   cfg.Path,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			LocalTime:  true,
			Compress:   true,
			Level:      logLevel,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat:  time.RFC3339,
				DisableTimestamp: false,
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyMsg:  "message",
					logrus.FieldKeyTime: "timestamp",
				},
			},
		}))
	} else {
		logrus.SetOutput(os.Stdout)

		if logLevel == logrus.DebugLevel {
			logrus.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC3339,
			})
		} else {
			logrus.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: time.RFC3339,
			})
		}
	}

	logrus.SetReportCaller(true)
}
