// Package handlers provides handlers for the HTTP API endpoints of the application.
package handlers

import (
	"net/http"

	"github.com/clementb49/welsh_academy/dto"
	"github.com/clementb49/welsh_academy/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler interface specifies methods for handling user-related operations
type UserHandler interface {
	CreateUserHandler(*gin.Context)
	LoginHandler(*gin.Context)
	GetUserByIdParamHandler(ctx *gin.Context)
	GetCurrentUserHandler(ctx *gin.Context)
	getUserByIdHandler(ctx *gin.Context, userId uint)
}

// userHandler implements UserHandler
type userHandler struct {
	service services.UserService
	logger  *zap.Logger
}

// NewUserHandler returns a new instance of userHandler
func NewUserHandler(service services.UserService) UserHandler {
	return &userHandler{
		service: service,
		logger:  zap.L(),
	}
}

// CreateUserHandler handles the creation of a new user
func (h *userHandler) CreateUserHandler(ctx *gin.Context) {
	var input dto.CreateUserReqBody
	err := ctx.ShouldBind(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.CreateUser(&input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

// LoginHandler handles user login
func (h *userHandler) LoginHandler(ctx *gin.Context) {
	var input dto.LoginReqBody
	err := ctx.ShouldBind(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.LoginUser(&input)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

// getUserByIdHandler handles getting a user by their ID
func (h *userHandler) getUserByIdHandler(ctx *gin.Context, userId uint) {
	user, err := h.service.GetUserById(userId)
	if err != nil {
		gormErrorResponseHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// GetUserByIdParamHandler handles getting a user by their ID specified in the path parameter
func (h *userHandler) GetUserByIdParamHandler(ctx *gin.Context) {
	var input dto.CommonIdPathUri
	err := ctx.ShouldBindUri(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.getUserByIdHandler(ctx, input.ID)
}

// GetCurrentUserHandler handles getting the current user
func (h *userHandler) GetCurrentUserHandler(ctx *gin.Context) {
	userId := ctx.GetUint("userId")
	h.getUserByIdHandler(ctx, userId)
}
