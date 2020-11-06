package router

import (
	"strconv"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/metric"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	labelEcoCode   = "code"
	labelEcoMethod = "method"
	labelEcoHost   = "host"
	labelEcoURL    = "url"
)

// Metrics keeps echo Prometheus metrics.
type Metrics struct {
	ReqQPS      *prometheus.CounterVec
	ReqDuration *prometheus.HistogramVec
}

//nolint:gochecknoglobals
var (
	metrics = Metrics{
		ReqQPS: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: metric.Namespace,
			Name:      "http_request_total",
			Help:      "The total http requests received",
		}, []string{labelEcoCode, labelEcoMethod, labelEcoHost, labelEcoURL}),

		ReqDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: metric.Namespace,
			Name:      "http_request_duration_seconds",
			Help:      "A histogram of latencies for requests received",
		}, []string{labelEcoCode, labelEcoMethod, labelEcoHost, labelEcoURL}),
	}
)

func prometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}

			uri := req.URL.Path
			status := strconv.Itoa(res.Status)
			duration := time.Since(start).Seconds()

			metrics.ReqQPS.WithLabelValues(status, req.Method, req.Host, uri).Inc()
			metrics.ReqDuration.WithLabelValues(status, req.Method, req.Host, uri).Observe(duration)

			return nil
		}
	}
}
