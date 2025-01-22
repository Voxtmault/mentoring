package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/voxtmault/mentoring/config"
	"github.com/voxtmault/mentoring/internal/routes"
	"github.com/voxtmault/mentoring/pkg/logger"
	"github.com/voxtmault/mentoring/pkg/storage/mariadb"
	"github.com/voxtmault/mentoring/pkg/storage/redis"
	val "github.com/voxtmault/mentoring/pkg/validator"
)

func main() {
	// Read / Load the configuration from .env file
	cfg := config.New(".env")

	// Initialize the logger
	if err := logger.InitLogger(&cfg.LoggingConfig); err != nil {
		panic(err)
	}

	// Open the initial connection to mariadb
	if err := mariadb.InitMariaDB(&cfg.DBConfig); err != nil {
		panic(err)
	}

	// Open the initial connection to redis
	if err := redis.InitRedis(&cfg.RedisConfig); err != nil {
		panic(err)
	}

	// Initialize the validator
	if err := val.InitValidator(); err != nil {
		panic(err)
	}

	// Init HTTP Routes
	echoInstance, err := routes.Init(cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		if cfg.SSLConfig.CertPath != "" && cfg.SSLConfig.KeyPath != "" {
			// Try HTTPS
			slog.Info("running restful service on https", "port", cfg.AppPort)
			echoInstance.Pre(middleware.HTTPSRedirect())

			s := http.Server{
				Addr:    fmt.Sprintf(":%s", cfg.AppPort),
				Handler: echoInstance,
			}
			if err := s.ListenAndServeTLS(cfg.SSLConfig.CertPath, cfg.SSLConfig.KeyPath); err != nil {
				slog.Warn("shutting down the service...")
			}
		} else {
			// Use HTTP as default if no ssl config is provided or valid
			slog.Info("running restful service on http", "port", cfg.AppPort)
			if err := echoInstance.Start(":" + cfg.AppPort); err != nil {
				slog.Warn("shutting down the service...")
			}
		}
	}()

	// ensure gracefull shutdown
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt, syscall.SIGTERM)
	<-interupt

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	slog.Warn("received shutdown signal. killing service...")

	// Close the main echo server
	if err := echoInstance.Shutdown(ctx); err != nil {
		slog.Error("error shutting down echo server", "reason", err)
	}

	// Close the database connection
	if err := mariadb.Close(); err != nil {
		slog.Error("error closing database connection", "reason", err)
	}
	if err := redis.CloseRedis(); err != nil {
		slog.Error("error closing redis connection", "reason", err)
	}

	// Finally close the logger
	if err := logger.CloseLogger(); err != nil {
		slog.Error("error closing logger", "reason", err)
	}
	slog.Info("closed all connections. service has gracefully shutdown. Bye :)")
}
