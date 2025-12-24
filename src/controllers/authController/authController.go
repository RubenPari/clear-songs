/**
 * Authentication Controller Package
 *
 * This package handles all authentication-related HTTP endpoints for the Clear Songs API.
 * It manages the Spotify OAuth 2.0 flow, including login initiation, callback handling,
 * session management, and authentication status checking.
 *
 * OAuth Flow:
 * 1. User requests /auth/login -> redirects to Spotify authorization page
 * 2. User authorizes app -> Spotify redirects to /auth/callback?code=XXX
 * 3. Backend exchanges code for access token -> stores in cache/session
 * 4. User is authenticated -> can access protected routes
 *
 * The controller uses offline access tokens (AccessTypeOffline) to allow API calls
 * without requiring user re-authentication, improving user experience.
 *
 * @package authController
 * @author Clear Songs Development Team
 */
package authController

import (
	"context"
	"log"
	"os"

	cacheManager "github.com/RubenPari/clear-songs/src/cache"

	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

/**
 * Login initiates the Spotify OAuth 2.0 authentication flow
 *
 * This handler redirects the user to Spotify's authorization page where they
 * can grant permissions to the application. The function uses offline access
 * type to obtain a refresh token, allowing the application to make API calls
 * in the future without requiring the user to authenticate again.
 *
 * Process:
 * 1. Generates OAuth authorization URL with offline access type
 * 2. Includes a state parameter for security (CSRF protection)
 * 3. Redirects user to Spotify authorization page
 *
 * After user authorization, Spotify redirects back to /auth/callback with
 * an authorization code that can be exchanged for an access token.
 *
 * @param c - Gin context containing HTTP request and response
 */
func Login(c *gin.Context) {
	url := utils.GetOAuth2Config().AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Default().Printf("Redirecting to %s", url)

	c.Redirect(302, url)
}

// Callback handles the callback received from Spotify after the user
// login request.
// Performs the following logic:
// 1. Gets the authorization code received in the query string.
// 2. Exchanges the code with the access token.
// 3. Creates a Spotify client using the access token.
// 4. Saves the Spotify client to the session.
// 5. Makes a test call to verify authentication.
// 6. Returns a success JSON if authentication is successful,
// otherwise returns an error JSON.
func Callback(c *gin.Context) {
	code := c.Query("code")

	token, errToken := utils.GetOAuth2Config().Exchange(context.Background(), code)

	if errToken != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Error authenticating user",
		})
		return
	}

	// set token in cache for session retrieval and in service instance
	if err := cacheManager.SetToken(token); err != nil {
		log.Default().Printf("ERROR: Failed to save token to cache: %v", err)
	} else {
		log.Default().Println("Token saved to cache successfully")
	}
	utils.SpotifySvc.SetAccessToken(token)

	log.Default().Println("Called callback, created spotify wrapper")

	// get user info for testing
	user, errUser := utils.SpotifySvc.GetSpotifyClient().CurrentUser()

	if errUser != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Error authenticating user",
		})
		return
	}

	// set spotifyService instance for dependency injection
	c.Set("spotifyService", utils.SpotifySvc)

	// Store user info in response for potential API usage
	// (though we redirect, this ensures user data is available if needed)
	_ = user

	// Redirect to frontend callback page after successful authentication
	// The token has already been exchanged and saved, so we redirect to the frontend
	// callback route which will verify authentication status and redirect to dashboard.
	// We don't pass the code because it's already been used to exchange for the token.
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4200"
	}
	c.Redirect(302, frontendURL+"/callback")
}

// Logout clears the user authentication by clearing the Spotify client
// from the session.
// Performs the following logic:
// 1. Clears the Spotify client from the session.
// 2. Returns a success JSON if authentication is successful,
// otherwise returns an error JSON.
func Logout(c *gin.Context) {
	utils.SpotifySvc.SetAccessToken(nil)
	// Clear token from cache by setting it to nil (this also clears in-memory token)
	_ = cacheManager.SetToken(nil)

	log.Default().Println("Called logout, deleted client from session")

	c.JSON(200, gin.H{
		"success": true,
		"message": "User logged out",
	})
}

// IsAuth checks if the user is authenticated.
// Performs the following logic:
// 1. Checks if Spotify client is set.
// 2. If Spotify client is not set, returns an error JSON with
// status 401 and an error message.
// 3. If Spotify client is set, returns a success JSON with
// status 200 and a success message.
func IsAuth(c *gin.Context) {
	client := utils.SpotifySvc.GetSpotifyClient()
	if client == nil {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	user, err := client.CurrentUser()
	if err != nil {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "User authenticated",
		"user": gin.H{
			"spotify_id":   user.ID,
			"display_name": user.DisplayName,
			"email":        user.Email,
			"profile_image": func() string {
				if len(user.Images) > 0 {
					return user.Images[0].URL
				}
				return ""
			}(),
		},
	})
}
