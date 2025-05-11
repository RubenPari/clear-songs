package middlewares

import (
	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := cacheManager.Get("spotify_token")
		if value != nil {
			token := value.(*oauth2.Token)
			utils.SpotifySvc.SetAccessToken(token)
			c.Set("spotifyService", utils.SpotifySvc)
		}
		c.Next()
	}
}
