package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/voxtmault/mentoring/internal/controllers"
)

func userRoute(root *echo.Group) {
	user := root.Group("/user")

	user.POST("", controllers.AddUserController)
	// user.GET("", GetUser)
	// user.PUT("", UpdateUser)
	// user.DELETE("", DeleteUser)
}
