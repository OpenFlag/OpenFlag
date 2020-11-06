package metric

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// Namespace is the Prometheus metric namespace variable.
	Namespace = "openflag"

	labelDbName    = "db_name"
	labelRedisName = "redis_name"
)

// Metrics keeps global Prometheus metrics.
type Metrics struct {
	DbConnectionStatus    *prometheus.GaugeVec
	RedisConnectionStatus *prometheus.GaugeVec
}

//nolint:gochecknoglobals
var (
	metrics = Metrics{
		DbConnectionStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "db_connection_status",
			Help:      "Database connection status",
		}, []string{labelDbName}),
		RedisConnectionStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "redis_connection_status",
			Help:      "Redis connection status",
		}, []string{labelRedisName}),
	}
)

// ReportDbStatus reports status of database connection to the Prometheus.
func ReportDbStatus(db *gorm.DB, dbName string) {
	// 1 means query is ok and 0 means query is not ok
	status := 1
	if err := db.Exec("SELECT 1;").Error; err != nil {
		status = 0
	}

	metrics.DbConnectionStatus.With(prometheus.Labels{labelDbName: dbName}).Set(float64(status))
}

// ReportRedisStatus reports status of redis connection to the Prometheus.
func ReportRedisStatus(cmdable redis.Cmdable, redisName string) {
	// 1 means ping is ok and 0 means ping is not ok
	status := 1
	if err := cmdable.Ping().Err(); err != nil {
		status = 0
	}

	metrics.RedisConnectionStatus.With(prometheus.Labels{labelRedisName: redisName}).Set(float64(status))
}
