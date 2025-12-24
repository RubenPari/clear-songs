/**
 * Routes Package
 * 
 * This package defines all HTTP routes for the Clear Songs backend API.
 * It sets up route groups, applies middleware, and connects routes to their
 * respective controller handlers.
 * 
 * Route Structure:
 * - /auth: Authentication endpoints (login, callback, logout, is-auth)
 * - /track: Track management endpoints (summary, delete by artist, delete by range)
 * - /playlist: Playlist management endpoints (delete tracks, delete tracks and library)
 * 
 * Middleware:
 * - SessionMiddleware: Manages user sessions and authentication state
 * - CacheInvalidationMiddleware: Handles cache invalidation for data consistency
 * - SpotifyAuthMiddleware: Validates Spotify OAuth tokens for protected routes
 * 
 * All routes return JSON responses following a consistent API response format.
 * 
 * @package routes
 * @author Clear Songs Development Team
 */
package routes

import (
	"github.com/RubenPari/clear-songs/src/controllers/authController"
	"github.com/RubenPari/clear-songs/src/controllers/playlistController"
	"github.com/RubenPari/clear-songs/src/controllers/trackController"
	"github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

/**
 * SetUpRoutes configures all HTTP routes for the application
 * 
 * This function:
 * 1. Applies global middleware (session, cache invalidation)
 * 2. Sets up route groups for different feature areas
 * 3. Connects routes to controller handlers
 * 4. Configures 404 handler for unknown routes
 * 
 * Route Groups:
 * - /auth: Public authentication routes (no auth required)
 * - /track: Protected track management routes (requires Spotify auth)
 * - /playlist: Protected playlist management routes (requires Spotify auth)
 * 
 * @param server - The Gin engine instance to configure routes on
 */
func SetUpRoutes(server *gin.Engine) {
	/**
	 * Global Middleware
	 * 
	 * These middleware functions are applied to all routes:
	 * - SessionMiddleware: Manages user sessions, validates cookies, handles authentication state
	 * - CacheInvalidationMiddleware: Invalidates cache when data is modified to ensure consistency
	 */
	server.Use(middlewares.SessionMiddleware())
	server.Use(middlewares.CacheInvalidationMiddleware())

	/**
	 * 404 Not Found Handler
	 * 
	 * Handles requests to routes that don't exist.
	 * Returns a JSON error response with 404 status code.
	 */
	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "not found path",
		})
	})

	/**
	 * Authentication Routes Group
	 * 
	 * Public routes for user authentication. These routes do not require
	 * authentication and handle the Spotify OAuth flow.
	 * 
	 * Routes:
	 * - GET /auth/login: Initiates Spotify OAuth flow, redirects to Spotify
	 * - GET /auth/callback: Handles OAuth callback, exchanges code for token
	 * - GET /auth/logout: Logs out user and destroys session
	 * - GET /auth/is-auth: Checks if user is currently authenticated
	 */
	auth := server.Group("/auth")
	{
		auth.GET("/login", authController.Login)
		auth.GET("/callback", authController.Callback)
		auth.GET("/logout", authController.Logout)
		auth.GET("/is-auth", authController.IsAuth)
	}

	/**
	 * Track Management Routes Group
	 * 
	 * Protected routes for managing tracks. All routes require Spotify authentication
	 * via SpotifyAuthMiddleware, which validates the user's OAuth token.
	 * 
	 * Routes:
	 * - GET /track/summary: Get summary of tracks grouped by artist (optional min/max filters)
	 * - DELETE /track/by-artist/:id_artist: Delete all tracks from a specific artist
	 * - DELETE /track/by-range: Delete tracks based on count range (query params: min, max)
	 */
	track := server.Group("/track")
	{
		track.GET("/summary",
			middlewares.SpotifyAuthMiddleware(),
			trackController.GetTrackSummary)
		track.DELETE("/by-artist/:id_artist",
			middlewares.SpotifyAuthMiddleware(),
			trackController.DeleteTrackByArtist)
		track.DELETE("/by-range",
			middlewares.SpotifyAuthMiddleware(),
			trackController.DeleteTrackByRange)
	}

	/**
	 * Playlist Management Routes Group
	 * 
	 * Protected routes for managing Spotify playlists. All routes require Spotify
	 * authentication and playlist ownership/editing permissions.
	 * 
	 * Routes:
	 * - DELETE /playlist/delete-tracks: Remove all tracks from a playlist (tracks remain in library)
	 * - DELETE /playlist/delete-tracks-and-library: Remove tracks from playlist AND user's library (with backup)
	 * 
	 * Both routes accept playlist ID as query parameter: ?id=PLAYLIST_ID
	 */
	playlist := server.Group("/playlist")
	{
		playlist.GET("/list",
			middlewares.SpotifyAuthMiddleware(),
			playlistController.GetUserPlaylists)
		playlist.DELETE("/delete-tracks",
			middlewares.SpotifyAuthMiddleware(),
			playlistController.DeleteAllPlaylistTracks)
		playlist.DELETE("/delete-tracks-and-library",
			middlewares.SpotifyAuthMiddleware(),
			playlistController.DeleteAllPlaylistAndUserTracks)
	}
}
