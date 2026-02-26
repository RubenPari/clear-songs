package middleware

import (
	"context"
	"log"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
)

// SessionMiddlewareRefactored creates a session middleware that uses dependency injection
func SessionMiddlewareRefactored(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// Retrieve OAuth token from cache
		token, err := cacheRepo.GetToken(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to retrieve token from cache: %v", err)
		}

		// If token exists, user is authenticated
		if token != nil {
			// Configure the Spotify repository with the user's token
			if err := spotifyRepo.SetAccessToken(token); err != nil {
				log.Printf("ERROR: Failed to set access token: %v", err)
			} else {
				// Store Spotify repository in context for use by handlers
				c.Set("spotifyRepository", spotifyRepo)
			}
		} else {
			// Log when token is not found (for debugging)
			// Only log for non-auth endpoints to avoid spam
			if c.Request.URL.Path != "/auth/is-auth" {
				log.Printf("DEBUG: No token found in cache for path: %s", c.Request.URL.Path)
			}
		}

		// Continue to next middleware or handler
		c.Next()
	}
}
