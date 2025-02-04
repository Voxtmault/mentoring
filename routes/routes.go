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

	// user
	userRoutes := api.Group("/users")
	userRoutes.Use(middlewares.CheckAPIKey)

	userRoutes.GET("", controller.GetAllUsers)
	userRoutes.GET("/:id", controller.GetUserByID)
	userRoutes.PUT("/:id", controller.UpdateUser)
	userRoutes.DELETE("/:id", controller.DeleteUser)

	// user status
	userRoutes.GET("/status", controller.GetAllUserStatus)
	userRoutes.GET("/status/:id", controller.GetUserStatusById)
	userRoutes.PUT("/status/:id", controller.UpdateUserStatus)
	userRoutes.DELETE("/status/:id", controller.DeleteUserStatusById)

	// user role
	userRoutes.GET("/role", controller.GetAllUserRole)
	userRoutes.GET("/role/:id", controller.GetUserRoleById)
	userRoutes.PUT("/role/:id", controller.UpdateUserRole)
	userRoutes.DELETE("/role/:id", controller.DeleteUserRoleById)

	// user transaction
	userRoutes.GET("/transactions/:id", controller.GetUserTransactions)

	// stock
	stockRoutes := api.Group("/stock")
	stockRoutes.Use(middlewares.CheckAPIKey)
	stockRoutes.GET("", controller.GetAllStockItem)
	stockRoutes.GET("/:id", controller.GetStockItem)
	stockRoutes.PUT("/:id", controller.UpdateStockItem)
	stockRoutes.POST("", controller.CreateStockItem)
	stockRoutes.DELETE("/:id", controller.DeleteStockItem)

	// supplier
	supplierRoutes := api.Group("/supplier")
	supplierRoutes.Use(middlewares.CheckAPIKey)
	supplierRoutes.GET("", controller.GetAllSupplier)
	supplierRoutes.GET("/:id", controller.GetSupplierById)
	supplierRoutes.PUT("/:id", controller.UpdateSupplier)
	supplierRoutes.POST("", controller.CreateSupplier)
	supplierRoutes.DELETE("/:id", controller.DeleteSupplier)

	// Transaction
	transactionRoutes := api.Group("/transaction")
	transactionRoutes.Use(middlewares.CheckAPIKey)
	transactionRoutes.GET("", controller.GetAllTransaction)
	transactionRoutes.GET("/:id", controller.GetTransactionById)
	transactionRoutes.GET("/detail/:id", controller.GetTransactionDetailedById)
	transactionRoutes.PUT("/:id", controller.UpdateTransaction)
	transactionRoutes.POST("", controller.CreateTransaction)
	transactionRoutes.DELETE("/:id", controller.DeleteTransactionById)

	// Transaction Status
	transactionRoutes.GET("/status", controller.GetAllTransactionStatus)
	transactionRoutes.GET("/status/:id", controller.GetTransactionStatusById)
	transactionRoutes.PUT("/status/:id", controller.UpdateTransactionStatus)
	transactionRoutes.POST("/status", controller.CreateTransactionStatus)
	transactionRoutes.DELETE("/status/:id", controller.DeleteTransactionStatusById)

	return e
}
