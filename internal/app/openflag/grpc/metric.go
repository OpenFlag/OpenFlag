package grpc

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/metric"
	prom "github.com/OpenFlag/OpenFlag/pkg/monitoring/prometheus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

// nolint:gochecknoglobals
var metrics = grpc_prometheus.NewServerMetrics(func(opts *prometheus.CounterOpts) {
	opts.Namespace = metric.Namespace
})

// nolint:gochecknoinits
func init() {
	metrics.EnableHandlingTimeHistogram(func(opts *prometheus.HistogramOpts) {
		opts.Namespace = metric.Namespace
		opts.Buckets = prom.HistogramBuckets
	})
	prometheus.MustRegister(metrics)
}
