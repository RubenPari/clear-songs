package middlewares

import (
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if spotify client is set
		// and token header is set in request
		if utils.SpotifyClient == nil {
			c.AbortWithStatusJSON(401, gin.H{
				"status":  "error",
				"message": "Unauthorized",
			})
		}
		c.Next()
	}
}
