// Package handlers provides handlers for the HTTP API endpoints of the application.
package handlers

import (
	"net/http"

	"github.com/clementb49/welsh_academy/dto"
	"github.com/clementb49/welsh_academy/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecipeHandler is the interface for recipe handlers.
type RecipeHandler interface {
	CreateRecipeHandler(*gin.Context)
	GetAllRecipeHandler(*gin.Context)
	GetRecipeByIdHandler(ctx *gin.Context)
	DeleteRecipeById(ctx *gin.Context)
	AddToFavRecipeHandler(*gin.Context)
	DeleteFavRecipeHandler(ctx *gin.Context)
	GetAllFavRecipeHandler(ctx *gin.Context)
}

// recipeHandler is the implementation of RecipeHandler.
type recipeHandler struct {
	service services.RecipeService
	logger  *zap.Logger
}

// NewRecipeHandler creates a new instance of RecipeHandler.
func NewRecipeHandler(service services.RecipeService) RecipeHandler {
	return &recipeHandler{
		service: service,
		logger:  zap.L(),
	}
}

// CreateRecipeHandler is the handler for creating a new recipe.
func (h *recipeHandler) CreateRecipeHandler(ctx *gin.Context) {
	var input dto.RecipeReqBody
	err := ctx.ShouldBind(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.GetUint("userId")
	recipe, err := h.service.CreateRecipe(&input, userId)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, recipe)
}

// GetAllRecipeHandler is the handler for getting all recipes.
func (h *recipeHandler) GetAllRecipeHandler(ctx *gin.Context) {
	var input dto.CommonQueryPage
	err := ctx.ShouldBindQuery(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.PageSize == 0 {
		input.PageSize = 10
	}
	pageRecipes, err := h.service.GetAllRecipes(&input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, pageRecipes)
}

// GetRecipeByIdHandler is the handler for getting a recipe by ID.
func (h *recipeHandler) GetRecipeByIdHandler(ctx *gin.Context) {
	var input dto.CommonIdPathUri
	err := ctx.ShouldBindUri(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe, err := h.service.GetRecipeById(&input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, recipe)
}

// DeleteRecipeById is the handler for deleting a recipe by ID.
func (h *recipeHandler) DeleteRecipeById(ctx *gin.Context) {
	var input dto.CommonIdPathUri
	err := ctx.ShouldBindUri(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteRecipeById(&input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// AddToFavRecipeHandler is the handler for adding a recipe to favorites.
func (h *recipeHandler) AddToFavRecipeHandler(ctx *gin.Context) {
	var input dto.CommonIdPathUri
	err := ctx.ShouldBindUri(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.GetUint("userId")
	recipeRes, err := h.service.AddToFavRecipe(userId, &input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, recipeRes)
}

// DeleteFavRecipeHandler is the handler for deleting a recipe from favorites.
func (h *recipeHandler) DeleteFavRecipeHandler(ctx *gin.Context) {
	var input dto.CommonIdPathUri
	err := ctx.ShouldBindUri(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.GetUint("userId")
	err = h.service.DeleteFavRecipe(userId, &input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// GetAllFavRecipeHandler is the handler for getting all favorites recipes.
func (h *recipeHandler) GetAllFavRecipeHandler(ctx *gin.Context) {
	var input dto.CommonQueryPage
	err := ctx.ShouldBindQuery(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.GetUint("userId")
	if input.PageSize == 0 {
		input.PageSize = 10
	}
	pageRecipes, err := h.service.GetAllFavRecipes(&input, userId)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, pageRecipes)
}
