package middlewares

import (
	"github.com/gin-gonic/gin"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "not found path",
		})
	}
}
