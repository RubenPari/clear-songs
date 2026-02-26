/**
 * Middlewares Package - Session Middleware
 * 
 * This middleware manages user sessions and Spotify authentication state
 * for all incoming HTTP requests. It ensures that authenticated users have
 * their Spotify service instance available in the request context.
 * 
 * How it works:
 * 1. Retrieves the user's OAuth token from cache (if authenticated)
 * 2. If token exists, configures the Spotify service with the token
 * 3. Stores the Spotify service instance in the Gin context
 * 4. Allows the request to proceed to the next handler
 * 
 * The Spotify service instance stored in the context can be retrieved by
 * controllers and other middleware using: c.Get("spotifyService")
 * 
 * This middleware is applied globally to all routes, ensuring consistent
 * session management throughout the application.
 * 
 * @package middleware
 * @author Clear Songs Development Team
 */
package middleware

import (
	"log"

	cacheManager "github.com/RubenPari/clear-songs/internal/infrastructure/persistence/redis"
	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/gin-gonic/gin"
)

/**
 * SessionMiddleware manages user sessions and Spotify authentication
 * 
 * This middleware function:
 * - Retrieves the user's OAuth token from the cache (Redis)
 * - If a token exists, configures the global Spotify service instance
 * - Makes the Spotify service available to all handlers via Gin context
 * - Allows unauthenticated requests to proceed (for public routes)
 * 
 * The middleware runs before every request, ensuring that authenticated
 * users have their Spotify client ready for API calls.
 * 
 * @returns gin.HandlerFunc - Middleware function that processes requests
 */
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve OAuth token from cache (Redis)
		// Returns nil if user is not authenticated
		token := cacheManager.GetToken()
		
		// If token exists, user is authenticated
		if token != nil {
			// Configure the global Spotify service with the user's token
			// This allows the service to make authenticated API calls
			utils.SpotifySvc.SetAccessToken(token)
			
			// Store Spotify service in context for use by handlers
			// Handlers can retrieve it with: c.Get("spotifyService")
			c.Set("spotifyService", utils.SpotifySvc)
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
