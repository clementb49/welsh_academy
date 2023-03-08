// Package routes provides the routing configuration for the application.
package routes

import (
	"github.com/clementb49/welsh_academy/handlers"
	"github.com/clementb49/welsh_academy/repositories"
	"github.com/clementb49/welsh_academy/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitRecipeRoute initializes the routes for user-related HTTP requests
func InitUserRoutes(db *gorm.DB, unAuthRouter, authRouter *gin.RouterGroup) {
	logger := zap.S()
	logger.Debug("Initializing User routes ...")
	// Create a new instance of the user repository with the provided database
	userRepository := repositories.NewUserRepository(db)

	// Create a new instance of the user service with the user repository instance
	userService := services.NewUserService(userRepository)

	// Create a new instance of the user handler with the user service instance
	userHandler := handlers.NewUserHandler(userService)

	// Register the routes for user authentication and authorization
	unAuthRouter.POST("/register", userHandler.CreateUserHandler)
	unAuthRouter.POST("/login", userHandler.LoginHandler)

	// Define the HTTP routes for authenticated users
	authRouter.GET("/users/my", userHandler.GetCurrentUserHandler)
	authRouter.GET("/users/:id", userHandler.GetUserByIdParamHandler)
}
