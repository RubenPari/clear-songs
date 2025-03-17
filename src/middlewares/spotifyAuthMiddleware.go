package middlewares

import (
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
)

func SpotifyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service := utils.GetSpotifyService(c)
		if service == nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		client := service.GetSpotifyClient()
		if client == nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid Spotify client"})
			return
		}

		c.Set("spotifyClient", client)
		c.Next()
	}
}
