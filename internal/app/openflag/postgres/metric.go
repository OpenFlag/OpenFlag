package postgres

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/metric"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	LabelConnStat = "status"
)

type Metrics struct {
	Gauge *prometheus.GaugeVec
}

//nolint:gochecknoglobals
var (
	metrics = Metrics{
		Gauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: metric.Namespace,
			Name:      "postgres_connection_stats",
		}, []string{LabelConnStat}),
	}
)

func (pm Metrics) report(db *gorm.DB) {
	// 1 means query is ok and 0 means query is not ok
	queryStatus := 1
	if err := db.Exec("SELECT 1;").Error; err != nil {
		queryStatus = 0
	}

	stats := db.DB().Stats()

	pm.Gauge.With(prometheus.Labels{LabelConnStat: "open"}).Set(float64(stats.OpenConnections))
	pm.Gauge.With(prometheus.Labels{LabelConnStat: "in_use"}).Set(float64(stats.InUse))
	pm.Gauge.With(prometheus.Labels{LabelConnStat: "idle"}).Set(float64(stats.Idle))
	pm.Gauge.With(prometheus.Labels{LabelConnStat: "query_status"}).Set(float64(queryStatus))
}
