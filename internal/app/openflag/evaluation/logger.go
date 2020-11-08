package evaluation

import (
	"encoding/json"
	"io"

	"github.com/google/martian/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	// LoggerConfig represents a struct for evaluation logger configurations.
	LoggerConfig struct {
		Enabled    bool   `mapstructure:"enabled"`
		Path       string `mapstructure:"path" validate:"required"`
		MaxSize    int    `mapstructure:"max-size" validate:"required"`
		MaxBackups int    `mapstructure:"max-backups" validate:"required"`
		MaxAge     int    `mapstructure:"max-age" validate:"required"`
	}

	// Logger represents a struct for evaluation logger.
	Logger struct {
		Config    LoggerConfig `json:"config"`
		logWriter io.Writer
	}
)

// NewLogger creates a new evaluation logger.
func NewLogger(cfg LoggerConfig) Logger {
	logger := Logger{
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
func (l Logger) Log(result interface{}) {
	if !l.Config.Enabled {
		return
	}

	output, err := json.Marshal(result)
	if err != nil {
		log.Errorf("failed to marshal evaluation result to json: %s", err.Error())
	}

	_, err = l.logWriter.Write(output)
	if err != nil {
		log.Errorf("failed to write evaluation result: %s", err.Error())
	}
}
