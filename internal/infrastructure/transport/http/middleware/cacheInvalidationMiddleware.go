package middleware

import (
	"strings"

	cacheManager "github.com/RubenPari/clear-songs/internal/infrastructure/persistence/redis"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// CacheInvalidationMiddleware automatically invalidates cache based on the endpoint called
func CacheInvalidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Execute the request first
		c.Next()

		// Only invalidate cache if the request was successful
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			path := c.Request.URL.Path
			method := c.Request.Method

			// Only invalidate on modification operations (DELETE, POST, PUT, PATCH)
			if method == "DELETE" || method == "POST" || method == "PUT" || method == "PATCH" {
				invalidateBasedOnEndpoint(c, path)
			}
		}
	}
}

func invalidateBasedOnEndpoint(c *gin.Context, path string) {
	switch {
	case strings.HasPrefix(path, "/track/"):
		// Any track operation affects user data
		cacheManager.InvalidateUserData()

	case strings.HasPrefix(path, "/playlist/"):
		// Playlist operations
		if playlistID := c.Query("id"); playlistID != "" {
			cacheManager.InvalidatePlaylist(spotifyAPI.ID(playlistID))
		}

		// If it's a playlist operation that also affects user library
		if strings.Contains(path, "all") || strings.Contains(path, "library") {
			cacheManager.InvalidateUserData()
		}

	case strings.HasPrefix(path, "/album/"):
		// Album operations usually affect user library
		cacheManager.InvalidateUserData()

	default:
		// For any other modification operation, do a full reset as a safety measure
		cacheManager.Reset()
	}
}
