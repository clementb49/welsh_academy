package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/clementb49/welsh_academy/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Defines a custom error type for unauthorized access
type UnathorizatedError struct {
	Status  string `json:"status"`  // HTTP status string
	Code    int    `json:"code"`    // HTTP status code
	Method  string `json:"method"`  // HTTP method
	Message string `json:"message"` // Error message
}

// Helper function to create a Forbidden error response
func forbiddenError(ctx *gin.Context) *UnathorizatedError {
	return &UnathorizatedError{
		Status:  "Forbidden",                                // Set the status string
		Code:    http.StatusForbidden,                       // Set the HTTP status code
		Method:  ctx.Request.Method,                         // Get the HTTP method used for the request
		Message: "Authorization required for this endpoint", // Set the error message
	}
}

// Helper function to create an Unauthorized error response
func unauthorizedError(ctx *gin.Context) *UnathorizatedError {
	return &UnathorizatedError{
		Status:  "Unauthorized",                    // Set the status string
		Code:    http.StatusUnauthorized,           // Set the HTTP status code
		Method:  ctx.Request.Method,                // Get the HTTP method used for the request
		Message: "Access token invalid or expired", // Set the error message
	}
}

// Auth returns a middleware function that validates a JWT token in the Authorization header.
func Auth() gin.HandlerFunc {
	// Create a logger instance.
	logger := zap.L()
	// Return a gin.HandlerFunc, which is a function that takes a Context instance as an argument.
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// Get the Authorization header value.
		authStr := ctx.GetHeader("Authorization")
		// Check if the Authorization header is missing.
		if authStr == "" {
			// Create an error response object for forbidden access.
			errorResponse := forbiddenError(ctx)
			logger.Sugar().Errorw("Authorization header invalid", errorResponse)
			// Send the error response as JSON and abort the request.
			ctx.JSON(http.StatusUnauthorized, errorResponse)
			defer ctx.AbortWithStatus(http.StatusForbidden)
		}
		// Split the Authorization header value into two parts: "Bearer " and the JWT token.
		authStrs := strings.SplitAfter(authStr, "Bearer")
		if len(authStrs) != 2 {
			// Create an error response object for malformed Authorization header.
			errorResponse := forbiddenError(ctx)
			logger.Sugar().Error("malformed authorization header", errorResponse)
			// Send the error response as JSON and abort the request.
			ctx.JSON(http.StatusUnauthorized, errorResponse)
			defer ctx.AbortWithStatus(http.StatusForbidden)
		}
		// Extract the token string by trimming whitespace from the second split element.
		tokenStr := strings.Trim(authStrs[1], " ")
		// Verify the token string using the VerifyToken function from a `utils` package. This function returns the token claims and an error (if any).
		claims, err := utils.VerifyToken(tokenStr)
		// If there is an error, return an "unauthorized" error response and abort the request. Otherwise, extract the user ID from the token claims and set it as a value on the `Context` object using `ctx.Set()`. Finally, call `ctx.Next()` to pass the request to the next middleware in the chain.
		if err != nil {
			errorResponse := unauthorizedError(ctx)
			logger.Sugar().Errorf("Error when validating jwt %w", err.Error())
			ctx.JSON(http.StatusUnauthorized, errorResponse)
			defer ctx.AbortWithStatus(http.StatusUnauthorized)
		} else {
			userId, err := strconv.Atoi(claims.ID)
			if err != nil {
				errorResponse := unauthorizedError(ctx)
				ctx.JSON(http.StatusUnauthorized, errorResponse)
				defer ctx.AbortWithStatus(http.StatusUnauthorized)
			}
			ctx.Set("userId", uint(userId))
			ctx.Next()
		}
	})
}
