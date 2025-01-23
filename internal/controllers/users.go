package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/voxtmault/mentoring/internal/models"
)

func AddUserController(c echo.Context) error {
	var payload models.User
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]interface{}{
			"message":       "unable to bind http request body",
			"error_message": err,
		})
	}

	return c.NoContent(http.StatusOK)
}
