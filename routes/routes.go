package routes

import (
	"tutor_rest/controller"
	"tutor_rest/middlewares"

	// "log"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()
	e.Use(middlewares.CORSMiddleware())

	// e.Static("/uploads", "uploads")

	api := e.Group("/api")

	api.POST("/users/login", controller.LoginUser)               
	api.POST("/users", controller.AddUser)                 

	userRoutes := api.Group("/users")
	userRoutes.Use(middlewares.CheckAPIKey)

	userRoutes.GET("", controller.GetAllUsers)       
	userRoutes.GET("/:id", controller.GetUserByID)    
	userRoutes.PUT("/:id", controller.UpdateUser)
	userRoutes.DELETE("/:id", controller.DeleteUser)             

	userRoutes.GET("/transactions", controller.GetUserTransactions) 

	userRoutes.GET("/suppliers", controller.ListSuppliers)          
	userRoutes.GET("/stock", controller.ViewStock)                        
	userRoutes.PUT("/stock/:id", controller.UpdateStock)                  
	userRoutes.POST("/transactions", controller.AddTransaction)           
	userRoutes.POST("/purchases", controller.AddPurchaseFromSupplier)     

	// Stock management routes
	userRoutes.POST("/stock", controller.CreateStockItem)       
	userRoutes.GET("/stock/:id", controller.GetStockItem)       
	userRoutes.PUT("/stock/:id", controller.UpdateStockItem)    
	userRoutes.DELETE("/stock/:id", controller.DeleteStockItem) 

	// Transaction routes
	userRoutes.POST("/transactions", controller.CreateTransaction)       
	userRoutes.GET("/transactions/:id", controller.GetTransaction)       
	userRoutes.PUT("/transactions/:id", controller.UpdateTransaction)    
	userRoutes.DELETE("/transactions/:id", controller.DeleteTransaction) 

	return e
}


