package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres driver should have blank import
)

const (
	healthCheckInterval = 1
	maxAttempts         = 60
)

type Config struct {
	Host               string        `mapstructure:"host" validate:"required"`
	Port               int           `mapstructure:"port" validate:"required"`
	Username           string        `mapstructure:"user" validate:"required"`
	Password           string        `mapstructure:"pass" validate:"required"`
	DBName             string        `mapstructure:"dbname" validate:"required"`
	ConnectTimeout     time.Duration `mapstructure:"connect-timeout" validate:"required"`
	ConnectionLifetime time.Duration `mapstructure:"connection-lifetime" validate:"required"`
	MaxOpenConnections int           `mapstructure:"max-open-connections"`
	MaxIdleConnections int           `mapstructure:"max-idle-connections"`
}

func Create(cfg Config) (*gorm.DB, error) {
	url := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s connect_timeout=%d sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password,
		int(cfg.ConnectTimeout.Seconds()),
	)

	pgDb, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	pgDb.DB().SetConnMaxLifetime(cfg.ConnectionLifetime)
	pgDb.DB().SetMaxOpenConns(cfg.MaxOpenConnections)
	pgDb.DB().SetMaxIdleConns(cfg.MaxIdleConnections)

	return pgDb, nil
}

func WithRetry(fn func(cfg Config) (*gorm.DB, error), cfg Config) *gorm.DB {
	for i := 0; i < maxAttempts; i++ {
		pgDb, err := fn(cfg)
		if err == nil {
			return pgDb
		}

		logrus.Errorf(
			"cannot connect to postgres. Waiting %d second. Error is: %s",
			healthCheckInterval, err.Error(),
		)

		time.Sleep(healthCheckInterval * time.Second)
	}

	logrus.Fatalf("could not connect to postgres after %d attempts", maxAttempts)

	return nil
}
