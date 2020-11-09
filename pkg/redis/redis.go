package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// Options represents a struct for creating Redis connection configurations.
type Options struct {
	Sentinel        bool          `mapstructure:"sentinel"`
	MasterName      string        `mapstructure:"master-name"`
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

// Create creates a Redis connection.
func Create(address string, options Options, master bool) (client redis.Cmdable, closeFunc func() error) {
	var result *redis.Client

	if master && options.Sentinel {
		result = redis.NewFailoverClient(
			&redis.FailoverOptions{
				SentinelAddrs:   []string{address},
				MasterName:      options.MasterName,
				PoolSize:        options.PoolSize,
				DialTimeout:     options.DialTimeout,
				ReadTimeout:     options.ReadTimeout,
				WriteTimeout:    options.WriteTimeout,
				PoolTimeout:     options.PoolTimeout,
				IdleTimeout:     options.IdleTimeout,
				MinIdleConns:    options.MinIdleConns,
				MaxRetries:      options.MaxRetries,
				MinRetryBackoff: options.MinRetryBackoff,
				MaxRetryBackoff: options.MaxRetryBackoff,
			},
		)
	} else {
		result = redis.NewClient(
			&redis.Options{
				Addr:            address,
				Password:        options.Password,
				PoolSize:        options.PoolSize,
				DialTimeout:     options.DialTimeout,
				ReadTimeout:     options.ReadTimeout,
				WriteTimeout:    options.WriteTimeout,
				PoolTimeout:     options.PoolTimeout,
				IdleTimeout:     options.IdleTimeout,
				MinIdleConns:    options.MinIdleConns,
				MaxRetries:      options.MaxRetries,
				MinRetryBackoff: options.MinRetryBackoff,
				MaxRetryBackoff: options.MaxRetryBackoff,
			},
		)
	}

	return result, result.Close
}
