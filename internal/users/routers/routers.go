package routers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rotisserie/eris"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
	"github.com/voxtmault/mentoring/library-project/pkg/validator"
	"gorm.io/gorm"
)

var (
	startTime time.Time
)

func Init(db *gorm.DB, cfg *config.AppConfig) (*echo.Echo, error) {

	// init echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &validator.CustomValidator{Validator: validator.Init()}

	// register middlewares
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
			consoleLogFormat := fmt.Sprintf("[%s] TTF: %v, Status: %v, Type: %v, URI: %v, IP: %v, Proto: %v, Device: %v\n",
				time.Now().Format("2006-01-02 15:04:05"), v.Latency, v.Status, v.Method, v.URIPath, v.RemoteIP, v.Protocol, v.UserAgent)

			fmt.Println(consoleLogFormat)

			return nil
		},
	}))

	startTime = time.Now()
	root := e.Group(cfg.AppRoot)

	root.GET("/health_check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "alive",
			"uptime": time.Since(startTime).String(),
		})
	})

	// Call the route handler
	if err := userRoutes(root, cfg, db); err != nil {
		return nil, eris.Wrap(err, "registering user routes")
	}

	return e, nil
}
