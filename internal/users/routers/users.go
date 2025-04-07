package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/voxtmault/mentoring/library-project/internal/users/controllers"
	"github.com/voxtmault/mentoring/library-project/internal/users/services"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
	"gorm.io/gorm"
)

func userRoutes(root *echo.Group, cfg *config.AppConfig, db *gorm.DB) error {
	// Init the controller and service
	userService := services.New(db, cfg)
	userController := controllers.NewUserController(userService, cfg)

	// Register the HTTP RESTful Route
	user := root.Group("/user")
	user.GET("", userController.GetUsers)
	user.GET("/:id", userController.GetUserByID)
	user.POST("", userController.CreateUser)

	return nil
}
