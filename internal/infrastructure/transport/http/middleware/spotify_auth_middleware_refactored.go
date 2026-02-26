package middleware

import (
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// SpotifyAuthMiddlewareRefactored creates an auth middleware that uses dependency injection
func SpotifyAuthMiddlewareRefactored() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Spotify repository from context (set by SessionMiddleware)
		repo, exists := c.Get("spotifyRepository")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		spotifyRepo, ok := repo.(shared.SpotifyRepository)
		if !ok || spotifyRepo == nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid Spotify repository"})
			return
		}

		// Get client from repository (for backward compatibility)
		// Note: In a fully refactored version, we might not need this
		// as controllers would use the repository directly
		if spotifyRepoImpl, ok := spotifyRepo.(interface{ GetClient() *spotifyAPI.Client }); ok {
			client := spotifyRepoImpl.GetClient()
			if client == nil {
				c.AbortWithStatusJSON(401, gin.H{"message": "Invalid Spotify client"})
				return
			}
			c.Set("spotifyClient", client)
		}

		// Also store the repository for direct use
		c.Set("spotifyRepository", spotifyRepo)
		c.Next()
	}
}
