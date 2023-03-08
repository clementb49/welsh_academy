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

// InitRecipeRoute initializes the routes for recipe-related HTTP requests
func InitRecipeRoute(db *gorm.DB, unauthRouter, authRouter *gin.RouterGroup) {
	logger := zap.S()
	logger.Debug("Initializing recipe routes ...")

	// Create a new recipe repository using the provided database instance
	recipeRepository := repositories.NewRecipeRepository(db)
	// Create a new recipe service using the recipe repository
	recipeService := services.NewRecipeService(recipeRepository)
	// Create a new recipe handler using the recipe service
	recipeHandler := handlers.NewRecipeHandler(recipeService)

	// Define the HTTP routes for authenticated users
	authRouter.POST("/recipes", recipeHandler.CreateRecipeHandler)
	authRouter.DELETE("/recipes/:id", recipeHandler.DeleteRecipeById)
	authRouter.PATCH("/recipes/:id/favorite", recipeHandler.AddToFavRecipeHandler)
	authRouter.DELETE("/recipes/:id/favorite", recipeHandler.DeleteFavRecipeHandler)
	authRouter.GET("recipes/favorites", recipeHandler.GetAllFavRecipeHandler)

	// Define the HTTP routes for unauthenticated users
	unauthRouter.GET("/recipes", recipeHandler.GetAllRecipeHandler)
	unauthRouter.GET("/recipes/:id", recipeHandler.GetRecipeByIdHandler)
}
