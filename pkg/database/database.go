package database

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/mysql"    // MySQL driver should have blank import
	_ "github.com/jinzhu/gorm/dialects/postgres" // PostgreSQL driver should have blank import
)

const (
	healthCheckInterval = 1
	maxAttempts         = 60
)

// Options represents a struct for creating database connection configurations.
type Options struct {
	ConnectionLifetime time.Duration `mapstructure:"connection-lifetime"`
	MaxOpenConnections int           `mapstructure:"max-open-connections"`
	MaxIdleConnections int           `mapstructure:"max-idle-connections"`
}

// Create creates a database connection.
func Create(driver string, connStr string, options Options) (*gorm.DB, error) {
	db, err := gorm.Open(driver, connStr)
	if err != nil {
		return nil, err
	}

	db.DB().SetConnMaxLifetime(options.ConnectionLifetime)
	db.DB().SetMaxOpenConns(options.MaxOpenConnections)
	db.DB().SetMaxIdleConns(options.MaxIdleConnections)

	return db, nil
}

// WithRetry provides functionality for having retry for connecting to database.
func WithRetry(
	fn func(driver string, connStr string, options Options,
	) (*gorm.DB, error), driver string, connStr string, options Options) *gorm.DB {
	for i := 0; i < maxAttempts; i++ {
		db, err := fn(driver, connStr, options)
		if err == nil {
			return db
		}

		logrus.Errorf(
			"cannot connect to database. Waiting %d second. Error is: %s",
			healthCheckInterval, err.Error(),
		)

		time.Sleep(healthCheckInterval * time.Second)
	}

	logrus.Fatalf("could not connect to database after %d attempts", maxAttempts)

	return nil
}
