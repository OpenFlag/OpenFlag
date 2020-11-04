package postgres

import (
	"fmt"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"

	"github.com/carlescere/scheduler"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres driver should have blank import
)

const (
	healthCheckInterval = 1
	maxAttempts         = 60
)

func Create(postgres config.Postgres) (*gorm.DB, error) {
	url := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s connect_timeout=%d sslmode=disable",
		postgres.Host, postgres.Port, postgres.Username, postgres.DBName, postgres.Password,
		int(postgres.ConnectTimeout.Seconds()),
	)

	postgresDb, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	postgresDb.DB().SetConnMaxLifetime(postgres.ConnectionLifetime)
	postgresDb.DB().SetMaxOpenConns(postgres.MaxOpenConnections)
	postgresDb.DB().SetMaxIdleConns(postgres.MaxIdleConnections)

	_, err = scheduler.Every(healthCheckInterval).Seconds().Run(func() { metrics.report(postgresDb) })
	if err != nil {
		return nil, err
	}

	return postgresDb, nil
}

func WithRetry(fn func(postgres config.Postgres) (*gorm.DB, error), postgres config.Postgres) *gorm.DB {
	for i := 0; i < maxAttempts; i++ {
		db, err := fn(postgres)
		if err == nil {
			return db
		}

		logrus.Errorf("postgres: cannot connect to postgres. waiting %d second. error is: %s",
			healthCheckInterval, err.Error())

		time.Sleep(healthCheckInterval * time.Second)
	}

	panic(fmt.Sprintf("postgres: could not connect to postgres after %d attempts", maxAttempts))
}
