package log

import (
	"io"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// RotateFileConfig represents a configuration struct for the Logrus file rotation hook.
type RotateFileConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LocalTime  bool
	Compress   bool
	Level      logrus.Level
	Formatter  logrus.Formatter
}

// RotateFileHook represents a struct for the Logrus file rotation hook.
type RotateFileHook struct {
	Config    RotateFileConfig
	logWriter io.Writer
}

// NewRotateFileHook created a new Logrus file rotation hook.
func NewRotateFileHook(config RotateFileConfig) logrus.Hook {
	hook := RotateFileHook{
		Config: config,
	}

	hook.logWriter = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		LocalTime:  config.LocalTime,
		Compress:   config.Compress,
	}

	return &hook
}

// Levels implements the Logrus Hook interface.
func (hook *RotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.Config.Level+1]
}

// Fire implements the Logrus Hook interface.
func (hook *RotateFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.logWriter.Write(b)
	if err != nil {
		return err
	}

	return nil
}
