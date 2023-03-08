// Package handlers provides handlers for the HTTP API endpoints of the application.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// gormErrorResponseHandler handles GORM errors and returns an appropriate HTTP status code along with an error message in JSON format.
// The error types are checked, and the corresponding HTTP status code is set. If the error is not recognized, a 500 Internal Server Error status code is returned.
// The function takes a gin.Context object and the GORM error as input parameters. The error message is included in the JSON response.
// The function is used in the Golang Gin framework to handle GORM database errors in HTTP request handlers.
func gormErrorResponseHandler(ctx *gin.Context, err error) {
	var httpStatus int
	switch err {
	case gorm.ErrRecordNotFound:
		httpStatus = http.StatusNotFound
	case gorm.ErrDuplicatedKey:
		httpStatus = http.StatusConflict
	default:
		httpStatus = http.StatusInternalServerError
	}
	ctx.JSON(httpStatus, gin.H{"error": err.Error()})
}
