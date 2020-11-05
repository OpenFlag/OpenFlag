package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type Config struct {
	Address         string        `mapstructure:"address" validate:"required"`
	Password        string        `mapstructure:"password"`
	PoolSize        int           `mapstructure:"pool-size"`
	MinIdleConns    int           `mapstructure:"min-idle-conns"`
	DialTimeout     time.Duration `mapstructure:"dial-timeout"`
	ReadTimeout     time.Duration `mapstructure:"read-timeout"`
	WriteTimeout    time.Duration `mapstructure:"write-timeout"`
	PoolTimeout     time.Duration `mapstructure:"pool-timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle-timeout"`
	MaxRetries      int           `mapstructure:"max-retries"`
	MinRetryBackoff time.Duration `mapstructure:"min-retry-backoff"`
	MaxRetryBackoff time.Duration `mapstructure:"max-retry-backoff"`
}

func Create(cfg Config) (client redis.Cmdable, closeFunc func() error) {
	result := redis.NewClient(
		&redis.Options{
			Addr:            cfg.Address,
			Password:        cfg.Password,
			PoolSize:        cfg.PoolSize,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			PoolTimeout:     cfg.PoolTimeout,
			IdleTimeout:     cfg.IdleTimeout,
			MinIdleConns:    cfg.MinIdleConns,
			MaxRetries:      cfg.MaxRetries,
			MinRetryBackoff: cfg.MinRetryBackoff,
			MaxRetryBackoff: cfg.MaxRetryBackoff,
		},
	)

	return result, result.Close
}
