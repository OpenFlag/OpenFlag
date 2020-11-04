package log

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type AppLogger struct {
	Level      string `mapstructure:"level" validate:"required"`
	Path       string `mapstructure:"path" validate:"required"`
	MaxSize    int    `mapstructure:"max-size" validate:"required"`
	MaxBackups int    `mapstructure:"max-backups" validate:"required"`
	MaxAge     int    `mapstructure:"max-age" validate:"required"`
	StdOut     bool   `mapstructure:"stdout"`
}

func SetupLogger(cfg AppLogger) {
	logLevel, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logLevel = logrus.DebugLevel
	}

	logrus.SetLevel(logLevel)

	if !cfg.StdOut {
		logrus.SetOutput(ioutil.Discard)

		rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   cfg.Path,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Level:      logLevel,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat:  time.RFC3339,
				DisableTimestamp: false,
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyMsg:  "message",
					logrus.FieldKeyTime: "timestamp",
				},
			},
		})
		if err != nil {
			logrus.Fatalf("failed to create log rotation: %s", err.Error())
		}

		logrus.AddHook(rotateFileHook)
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
