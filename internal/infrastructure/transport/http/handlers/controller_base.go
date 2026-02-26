package handlers

import (
	"github.com/RubenPari/clear-songs/internal/application/shared/dto"
	"github.com/gin-gonic/gin"
)

// BaseController provides common functionality for all controllers
type BaseController struct{}

// JSONSuccess sends a successful JSON response
func (bc *BaseController) JSONSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, dto.Success(data))
}

// JSONError sends an error JSON response
func (bc *BaseController) JSONError(c *gin.Context, status int, code, message string) {
	c.JSON(status, dto.ErrorResponse(code, message))
}

// JSONValidationError sends a validation error response
func (bc *BaseController) JSONValidationError(c *gin.Context, message string) {
	c.JSON(400, dto.ValidationError(message))
}

// JSONInternalError sends an internal server error response
func (bc *BaseController) JSONInternalError(c *gin.Context, message string) {
	c.JSON(500, dto.InternalError(message))
}

// JSONNotFound sends a not found error response
func (bc *BaseController) JSONNotFound(c *gin.Context, resource string) {
	c.JSON(404, dto.NotFoundError(resource))
}

// JSONUnauthorized sends an unauthorized error response
func (bc *BaseController) JSONUnauthorized(c *gin.Context) {
	c.JSON(401, dto.UnauthorizedError())
}
