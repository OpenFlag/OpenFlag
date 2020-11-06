package log

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/natefinch/lumberjack.v2"
)

// AccessLogger represents a struct for echo access logger middleware configurations.
type AccessLogger struct {
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path" validate:"required"`
	Format     string `mapstructure:"format"`
	MaxSize    int    `mapstructure:"max-size" validate:"required"`
	MaxBackups int    `mapstructure:"max-backups" validate:"required"`
	MaxAge     int    `mapstructure:"max-age" validate:"required"`
}

// LoggerMiddleware is an echo middleware for access logging.
func LoggerMiddleware(cfg AccessLogger) echo.MiddlewareFunc {
	if !cfg.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}
		}
	}

	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: &lumberjack.Logger{
			Filename:   cfg.Path,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   true,
		},
		Format: cfg.Format,
	})
}
