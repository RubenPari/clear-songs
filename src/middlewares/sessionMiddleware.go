package middlewares

import (
	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := cacheManager.GetToken()
		if token != nil {
			utils.SpotifySvc.SetAccessToken(token)
			c.Set("spotifyService", utils.SpotifySvc)
		}
		c.Next()
	}
}
