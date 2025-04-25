package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/voxtmault/mentoring/library-project/internal/users/routers"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
	"github.com/voxtmault/mentoring/library-project/pkg/logging"
	"github.com/voxtmault/mentoring/library-project/pkg/storage/mariadb"
)

func main() {
	slog.Info("starting users service")

	cwd, _ := os.Getwd()
	log.Printf("Current working directory: %s", cwd)

	slog.Info("reading provided env file")
	cfg := config.New("envs/users/.env")

	slog.Info("initializing mariadb connection and ORM instance")
	db, err := mariadb.Init(&cfg.MariaDBConfig)
	if err != nil {
		log.Fatalf("failed to init mariadb: %v", err)
	}
	orm, err := mariadb.InitORM(db)
	if err != nil {
		log.Fatalf("failed to init mariadb ORM: %v", err)
	}

	slog.Info("initializing logger")
	serLog, errLog := logging.InitLogger(&cfg.LoggingConfig)

	slog.Info("initializing echo instance")
	e, err := routers.Init(orm, cfg)
	if err != nil {
		log.Fatalf("failed to init echo: %v", err)
	}

	slog.Info("starting echo server")
	go func() {
		if cfg.SSLConfig.CertPath != "" && cfg.SSLConfig.KeyPath != "" {
			slog.Info("starting echo server with SSL", "port", cfg.AppPort)
			e.Pre(middleware.HTTPSRedirect())

			s := http.Server{
				Addr:    ":" + cfg.AppPort,
				Handler: e,
			}
			if err := s.ListenAndServeTLS(cfg.SSLConfig.CertPath, cfg.SSLConfig.KeyPath); err != nil {
				slog.Warn("failed to start echo server with SSL", "error", err)
				return
			}
		} else {
			slog.Warn("starting echo server without SSL", "port", cfg.AppPort)
			if err := e.Start(":" + cfg.AppPort); err != nil {
				slog.Warn("failed to start echo server without SSL", "error", err)
				return
			}
		}
	}()

	// ensure gracefull shutdown
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt, syscall.SIGTERM)
	<-interupt

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	slog.Info("gracefully shutting down echo server")
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("error shutting down echo server", "reason", err)
	}

	slog.Info("closing mariadb connection")
	if err := mariadb.Close(db); err != nil {
		slog.Error("error closing mariadb connection", "reason", err)
	}

	slog.Info("closing logger")
	if err := logging.CloseLogger(serLog, errLog); err != nil {
		slog.Error("error closing logger", "reason", err)
	}

	slog.Info("service stopped")
}
