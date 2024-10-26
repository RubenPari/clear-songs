package middlewares

import (
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if spotify client is set
		if _, err := utils.SpotifySvc.GetSpotifyClient().CurrentUser(); err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"status":  "error",
				"message": "Unauthorized",
			})
		} else {
			c.Next()
		}
	}
}
