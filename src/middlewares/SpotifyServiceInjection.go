package middlewares

import (
	"github.com/RubenPari/clear-songs/src/services/SpotifyService"
	"github.com/gin-gonic/gin"
)

// SpotifyServiceInjection returns a middleware that injects the given SpotifyService
// into the context under the key "spotifyService". This allows the service to be
// used in any handlers that are called after this middleware.
func SpotifyServiceInjection(service *SpotifyService.SpotifyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("spotifyService", service)
		c.Next()
	}
}
