package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/voxtmault/mentoring/library-project/internal/users/interfaces"
	"github.com/voxtmault/mentoring/library-project/internal/users/models"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
	"github.com/voxtmault/mentoring/library-project/pkg/http_utility"
)

type UserController struct {
	userService interfaces.User
	cfg         *config.AppConfig
}

func NewUserController(userService interfaces.User, cfg *config.AppConfig) *UserController {
	return &UserController{
		userService: userService,
		cfg:         cfg,
	}
}

func (uc *UserController) GetUsers(c echo.Context) error {
	var filter models.UserFilter
	res := http_utility.New(uc.cfg)
	if err := c.Bind(&filter); err != nil {
		res.StatusCode = http.StatusUnprocessableEntity
		res.InternalErrorMessage = http_utility.Unprocessable
		res.ErrorStack = eris.Wrap(err, "failed to bind user filters")
		res.Translate()

		return echo.NewHTTPError(res.StatusCode, res)
	}

	if err := c.Validate(filter); err != nil {
		res.StatusCode = http.StatusBadRequest
		res.InternalErrorMessage = http_utility.BadRequest
		res.ErrorStack = eris.Wrap(err, "failed to validate user filters")
		res.Translate()

		return echo.NewHTTPError(res.StatusCode, res)
	}

	res, _ = uc.userService.GetUsers(c.Request().Context(), &filter)
	if res.ErrorStack != nil {
		res.Translate()

		return echo.NewHTTPError(res.StatusCode, res)
	}

	return c.JSON(res.StatusCode, res)
}

func (uc *UserController) GetUserByID(c echo.Context) error {
	res := http_utility.New(uc.cfg)
	id := c.Param("id")
	convId, err := strconv.Atoi(id)
	if err != nil {
		res.StatusCode = http.StatusBadRequest
		res.InternalErrorMessage = http_utility.BadRequest
		res.ErrorStack = eris.Wrap(err, "invalid id value")
		res.Translate()

		return echo.NewHTTPError(res.StatusCode, res)
	}

	res, _ = uc.userService.GetUserByID(c.Request().Context(), uint(convId))
	if res.ErrorStack != nil {
		res.Translate()

		return echo.NewHTTPError(res.StatusCode, res)
	}

	return c.JSON(res.StatusCode, res)
}

func (uc *UserController) CreateUser(c echo.Context) error {
	var data models.User
	res := http_utility.New(uc.cfg)
	if err := c.Bind(&data); err != nil {
		res.StatusCode = http.StatusUnprocessableEntity
		res.InternalErrorMessage = http_utility.Unprocessable
		res.ErrorStack = eris.Wrap(err, "failed to bind user data")
		res.Translate()

		return echo.NewHTTPError(res.StatusCode, res)
	}

	if err := c.Validate(data); err != nil {
		res.StatusCode = http.StatusBadRequest
		res.InternalErrorMessage = http_utility.BadRequest
		res.ErrorStack = eris.Wrap(err, "failed to validate user data")
		res.Translate()

		return echo.NewHTTPError(res.StatusCode, res)
	}

	res, _ = uc.userService.CreateUser(c.Request().Context(), &data)
	if res.ErrorStack != nil {
		res.Translate()

		return echo.NewHTTPError(res.StatusCode, res)
	}

	return c.JSON(res.StatusCode, res)
}
