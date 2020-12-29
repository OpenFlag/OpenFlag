package engine

import (
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/metric"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	prom "github.com/OpenFlag/OpenFlag/pkg/monitoring/prometheus"
)

const (
	labelMethod        = "method"
	errorIncrementStep = 1
)

// Metrics keeps global Prometheus metrics.
type Metrics struct {
	ErrCounter *prometheus.CounterVec
	Histogram  *prometheus.HistogramVec
}

// nolint:gochecknoglobals
var (
	metrics = Metrics{
		ErrCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metric.Namespace,
				Name:      "engine_error_total",
			}, []string{labelMethod},
		),

		Histogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: metric.Namespace,
				Name:      "engine_duration_total",
				Buckets:   prom.HistogramBuckets,
			}, []string{labelMethod},
		),
	}
)

func (m Metrics) report(method string, startTime time.Time, err error) {
	if err != nil {
		m.ErrCounter.With(prometheus.Labels{labelMethod: method}).Add(errorIncrementStep)

		return
	}

	m.Histogram.With(prometheus.Labels{labelMethod: method}).
		Observe(time.Since(startTime).Seconds())
}
