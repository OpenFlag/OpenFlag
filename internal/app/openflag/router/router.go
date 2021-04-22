package router

import (
	"net/http"
	_ "net/http/pprof" // Golang pprof

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/OpenFlag/OpenFlag/pkg/log"
	"github.com/OpenFlag/OpenFlag/pkg/version"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// New creates a new application router.
func New(cfg config.Config) *echo.Echo {
	e := echo.New()

	debug := logrus.IsLevelEnabled(logrus.DebugLevel)

	e.Debug = debug

	e.HideBanner = true

	if !debug {
		e.HidePort = true
	}

	if debug {
		e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))
	}

	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Server.WriteTimeout = cfg.Server.WriteTimeout

	recoverConfig := middleware.DefaultRecoverConfig
	recoverConfig.DisablePrintStack = !debug
	e.Use(middleware.RecoverWithConfig(recoverConfig))

	e.Use(middleware.CORS())
	e.Use(version.Middleware)
	e.Use(log.LoggerMiddleware(cfg.Logger.AccessLogger))
	e.Use(prometheusMiddleware())

	return e
}
