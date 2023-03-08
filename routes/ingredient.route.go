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

// InitIngredientRoute initializes the routes for ingredient-related HTTP requests
func InitIngredientRoute(db *gorm.DB, unAuthRouter, authRouter *gin.RouterGroup) {
	logger := zap.S()
	logger.Debug("Initializing ingredient routes ...")
	// Create a new ingredient repository using the provided database instance
	ingredientRepository := repositories.NewIngredientRepository(db)
	// Create a new ingredient service using the ingredient repository
	ingredientService := services.NewIngredientService(ingredientRepository)
	// Create a new ingredient handler using the ingredient service
	ingredientHandler := handlers.NewIngredientHandlers(ingredientService)

	// Define the HTTP routes for authenticated users
	authRouter.POST("/ingredients", ingredientHandler.CreateIngredientHandler)
	authRouter.DELETE("ingredients/:id", ingredientHandler.DeleteIngredientByIdHandler)

	// Define the HTTP routes for unauthenticated users
	unAuthRouter.GET("/ingredients", ingredientHandler.GetAllIngredients)
	unAuthRouter.GET("/ingredients/:id", ingredientHandler.GetIngredientByIdHandler)
}
