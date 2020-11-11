package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Config represents a struct for starting the Prometheus monitoring server configurations.
type Config struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"address"`
}

// StartServer starts the Prometheus metric server for scraping the metrics data.
func StartServer(cfg Config) {
	if cfg.Enabled {
		metricServer := http.NewServeMux()
		metricServer.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(cfg.Address, metricServer); err != nil {
			logrus.Panicf("failed to start prometheus metrics server %s", err.Error())
		}
	}
}
