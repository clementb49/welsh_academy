// Package handlers provides handlers for the HTTP API endpoints of the application.
package handlers

import (
	"net/http"

	"github.com/clementb49/welsh_academy/dto"
	"github.com/clementb49/welsh_academy/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// IngredientHandlers defines the interface for ingredient handlers.
type IngredientHandlers interface {
	CreateIngredientHandler(ctx *gin.Context)
	GetAllIngredients(ctx *gin.Context)
	GetIngredientByIdHandler(*gin.Context)
	DeleteIngredientByIdHandler(ctx *gin.Context)
}

// ingredientHandlers is the implementation of the IngredientHandlers interface.
type ingredientHandlers struct {
	service services.IngredientService
	logger  *zap.Logger
}

// NewIngredientHandlers returns a new instance of IngredientHandlers.
func NewIngredientHandlers(service services.IngredientService) IngredientHandlers {
	return &ingredientHandlers{
		service: service,
		logger:  zap.L(),
	}
}

// CreateIngredientHandler creates a new ingredient with the given input.
func (h *ingredientHandlers) CreateIngredientHandler(ctx *gin.Context) {
	var input dto.IngredientReqBody
	err := ctx.ShouldBind(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ingredient, err := h.service.CreateIngredient(&input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, ingredient)
}

// GetAllIngredients returns a page of ingredients with the given input parameters.
func (h *ingredientHandlers) GetAllIngredients(ctx *gin.Context) {
	var input dto.CommonQueryPage
	err := ctx.ShouldBindQuery(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.PageSize == 0 {
		input.PageSize = 10
	}
	pageIngredients, err := h.service.GetAllIngredients(&input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, pageIngredients)
}

// GetIngredientByIdHandler returns the ingredient with the given ID.
func (h *ingredientHandlers) GetIngredientByIdHandler(ctx *gin.Context) {
	var input dto.CommonIdPathUri
	err := ctx.ShouldBindUri(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ingredient, err := h.service.GetIngredientById(input.ID)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, ingredient)
}

// DeleteIngredientByIdHandler deletes the ingredient with the given ID.
func (h *ingredientHandlers) DeleteIngredientByIdHandler(ctx *gin.Context) {
	var input dto.CommonIdPathUri
	err := ctx.ShouldBindUri(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteIngredientById(input.ID)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
	}
	ctx.Status(http.StatusNoContent)
}
