package engine

import (
	"encoding/json"
	"io"

	"github.com/sirupsen/logrus"

	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	// LoggerConfig represents a struct for evaluation logger configurations.
	LoggerConfig struct {
		Enabled    bool   `mapstructure:"enabled"`
		Path       string `mapstructure:"path"`
		MaxSize    int    `mapstructure:"max-size"`
		MaxBackups int    `mapstructure:"max-backups"`
		MaxAge     int    `mapstructure:"max-age"`
	}

	// EvaluationLogger represents a struct for evaluation logger.
	EvaluationLogger struct {
		Config    LoggerConfig
		logWriter io.Writer
	}
)

// Logger represents an evaluation result logger interface.
type Logger interface {
	Log(result Result)
}

// NewLogger creates a new evaluation logger.
func NewLogger(cfg LoggerConfig) EvaluationLogger {
	logger := EvaluationLogger{
		Config: cfg,
	}

	logger.logWriter = &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		LocalTime:  true,
		Compress:   true,
	}

	return logger
}

// Log logs the evaluation result.
func (l EvaluationLogger) Log(result Result) {
	if !l.Config.Enabled {
		return
	}

	output, err := json.Marshal(&result)
	if err != nil {
		logrus.Errorf("failed to marshal evaluation result to json: %s", err.Error())
	}

	output = []byte(string(output) + "\n")

	_, err = l.logWriter.Write(output)
	if err != nil {
		logrus.Errorf("failed to write evaluation result: %s", err.Error())
	}
}
