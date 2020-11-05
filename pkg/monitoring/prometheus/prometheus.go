package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"address" validate:"required"`
}

func StartPrometheusServer(cfg Config) {
	if cfg.Enabled {
		metricServer := http.NewServeMux()
		metricServer.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(cfg.Address, metricServer); err != nil {
			logrus.Panicf("failed to start prometheus metrics server %s", err.Error())
		}
	}
}
