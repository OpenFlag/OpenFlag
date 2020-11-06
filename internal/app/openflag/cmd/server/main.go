package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/metric"
	"github.com/OpenFlag/OpenFlag/pkg/monitoring/prometheus"
	"github.com/carlescere/scheduler"

	"github.com/OpenFlag/OpenFlag/pkg/redis"

	"github.com/OpenFlag/OpenFlag/pkg/postgres"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/router"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/spf13/cobra"
)

const (
	healthCheckInterval = 1
)

//nolint:funlen
func main(cfg config.Config) {
	e := router.New(cfg)

	pgDbMaster := postgres.WithRetry(postgres.Create, cfg.Postgres.Master)
	pgDbSlave := postgres.WithRetry(postgres.Create, cfg.Postgres.Slave)

	defer func() {
		if err := pgDbMaster.Close(); err != nil {
			logrus.Errorf("postgres master connection close error: %s", err.Error())
		}

		if err := pgDbSlave.Close(); err != nil {
			logrus.Errorf("postgres slave connection close error: %s", err.Error())
		}
	}()

	redisMasterClient, redisMasterClose := redis.Create(cfg.Redis.Master)
	redisSlaveClient, redisSlaveClose := redis.Create(cfg.Redis.Slave)

	defer func() {
		if err := redisMasterClose(); err != nil {
			logrus.Errorf("redis master connection close error: %s", err.Error())
		}

		if err := redisSlaveClose(); err != nil {
			logrus.Errorf("redis slave connection close error: %s", err.Error())
		}
	}()

	_, err := scheduler.Every(healthCheckInterval).Seconds().Run(func() {
		metric.ReportDbStatus(pgDbMaster, "pg_master")
		metric.ReportDbStatus(pgDbMaster, "pg_slave")
		metric.ReportRedisStatus(redisMasterClient, "redis_master")
		metric.ReportRedisStatus(redisSlaveClient, "redis_slave")
	})
	if err != nil {
		logrus.Fatalf("failed to start metric scheduler: %s", err.Error())
	}

	e.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	// ==========
	// Codes
	// ==========

	e.Static("/", "browser/openflag-ui/build")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := e.Start(cfg.Server.Address); err != nil {
			logrus.Fatalf("failed to start openflag server: %s", err.Error())
		}
	}()

	go prometheus.StartServer(cfg.Monitoring.Prometheus)

	logrus.Info("start openflag server!")

	s := <-sig

	logrus.Infof("signal %s received", s)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	e.Server.SetKeepAlivesEnabled(false)

	if err := e.Shutdown(ctx); err != nil {
		logrus.Errorf("failed to shutdown openflag server: %s", err.Error())
	}
}

// Register register server command for openflag binary.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run OpenFlag server component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
