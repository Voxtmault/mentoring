package middlewares

import (
    "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

// CORS Middleware Configuration
func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.DefaultCORSConfig)
}

func CheckAPIKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("X-API-KEY")
		if apiKey == "" {
			return echo.ErrUnauthorized
		}
		if apiKey != os.Getenv("API_KEY") {
			return echo.ErrForbidden
		}
		return next(c)
	}
}