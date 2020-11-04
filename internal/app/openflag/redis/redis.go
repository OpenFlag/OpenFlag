package redis

import (
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/go-redis/redis"
)

func Create(cfg config.RedisConfig) (client redis.Cmdable, closeFunc func() error) {
	result := redis.NewClient(
		&redis.Options{
			Addr:            cfg.Address,
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
