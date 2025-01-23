package routes

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/voxtmault/mentoring/config"

	"github.com/voxtmault/mentoring/pkg/logger"
	val "github.com/voxtmault/mentoring/pkg/validator"
)

var (
	startTime time.Time
)

func Init(cfg *config.AppConfig) (*echo.Echo, error) {
	// Initialize echo instance
	e := echo.New()
	e.HideBanner = true

	// Register Validator
	e.Validator = val.GetValidatorInstance()

	// Register Middlewares
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
		Level: 6,
	}))
	e.Use(middleware.Decompress())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogURIPath:   true,
		LogHost:      true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogProtocol:  true,
		LogUserAgent: true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			msg, err := json.Marshal(v)
			if err != nil {
				slog.Error("failed to marshal request log", "reason", err)
				return nil
			}

			// Suppress the console except for errors
			slog.Info("Request Log", "latency", v.Latency, "status", v.Status, "method", v.Method, "uri", v.URIPath, "ip", v.RemoteIP,
				"proto", v.Protocol, "device", v.UserAgent)

			if err := logger.LogRequest(msg); err != nil {
				slog.Error("failed to log request", "reason", err)
				return nil
			}

			return nil
		},
	}))
	e.Use(echoprometheus.NewMiddleware("fruit_store"))

	// /api/v1
	root := e.Group(cfg.AppRoot)
	startTime = time.Now()

	// Monitoring Endpoints
	root.GET("/health_check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "alive",
			"uptime": time.Since(startTime).String(),
		})
	})
	root.GET("/metrics", echoprometheus.NewHandler())

	// Register business routes
	userRoute(root)

	return e, nil
}
