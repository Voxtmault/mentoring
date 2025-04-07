package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/rotisserie/eris"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
)

func validateConfig(cfg *config.RedisConfig) error {
	if cfg.RedisHost == "" {
		return eris.New("redis host is empty")
	}
	if cfg.RedisPort == "" {
		return eris.New("redis port is empty")
	}
	if cfg.RedisPassword == "" {
		return eris.New("redis password is empty")
	}

	return nil
}

func Init(cfg *config.RedisConfig) (*redis.Client, error) {
	slog.Debug("init redis connection")
	if err := validateConfig(cfg); err != nil {
		return nil, eris.Wrap(err, "invalid redis configuration")
	}

	con := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       int(cfg.RedisDBNum),
	})

	if _, err := con.Ping(context.Background()).Result(); err != nil {
		return nil, eris.Wrap(err, "Init Redis")
	}

	slog.Info("redis connection established", "host", cfg.RedisHost)
	return con, nil
}

func Close(con *redis.Client) error {
	if err := con.Close(); err != nil {
		return eris.Wrap(err, "Closing redis connection")
	}

	return nil
}
