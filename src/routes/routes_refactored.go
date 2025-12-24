package routes

import (
	"github.com/RubenPari/clear-songs/src/infrastructure/di"
	"github.com/RubenPari/clear-songs/src/presentation/controllers"
	"github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

/**
 * SetUpRoutesRefactored configures all HTTP routes using dependency injection
 * 
 * This version uses the DI container to inject dependencies into controllers
 * and middleware, eliminating the need for global variables.
 * 
 * @param server - The Gin engine instance to configure routes on
 * @param container - The dependency injection container
 */
func SetUpRoutesRefactored(server *gin.Engine, container *di.Container) {
	/**
	 * Global Middleware
	 * 
	 * These middleware functions are applied to all routes:
	 * - SessionMiddlewareRefactored: Manages user sessions using DI
	 * - CacheInvalidationMiddleware: Invalidates cache when data is modified
	 */
	server.Use(middlewares.SessionMiddlewareRefactored(
		container.SpotifyRepo,
		container.CacheRepo,
	))
	server.Use(middlewares.CacheInvalidationMiddleware())

	/**
	 * 404 Not Found Handler
	 */
	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "not found path",
		})
	})

	/**
	 * Track Management Routes Group
	 * 
	 * Uses the complete refactored controller with dependency injection
	 */
	trackController := controllers.NewTrackControllerComplete(
		container.GetTrackSummaryUseCase,
		container.DeleteTracksByArtistUC,
		container.DeleteTracksByRangeUC,
		container.GetTracksByArtistUC,
	)
	
	track := server.Group("/track")
	{
		track.GET("/summary",
			middlewares.SpotifyAuthMiddlewareRefactored(),
			trackController.GetTrackSummary)
		track.GET("/by-artist/:id_artist",
			middlewares.SpotifyAuthMiddlewareRefactored(),
			trackController.GetTracksByArtist)
		track.DELETE("/by-artist/:id_artist",
			middlewares.SpotifyAuthMiddlewareRefactored(),
			trackController.DeleteTrackByArtist)
		track.DELETE("/by-range",
			middlewares.SpotifyAuthMiddlewareRefactored(),
			trackController.DeleteTrackByRange)
	}

	/**
	 * Authentication Routes Group
	 * 
	 * Uses the refactored controller with dependency injection
	 */
	authController := controllers.NewAuthControllerRefactored(
		container.LoginUC,
		container.CallbackUC,
		container.LogoutUC,
		container.IsAuthUC,
	)
	
	auth := server.Group("/auth")
	{
		auth.GET("/login", authController.Login)
		auth.GET("/callback", authController.Callback)
		auth.GET("/logout", authController.Logout)
		auth.GET("/is-auth", authController.IsAuth)
	}

	/**
	 * Playlist Management Routes Group
	 * 
	 * Uses the refactored controller with dependency injection
	 */
	playlistController := controllers.NewPlaylistControllerRefactored(
		container.GetUserPlaylistsUC,
		container.DeletePlaylistTracksUC,
		container.DeletePlaylistAndLibraryUC,
	)
	
	playlist := server.Group("/playlist")
	{
		playlist.GET("/list",
			middlewares.SpotifyAuthMiddlewareRefactored(),
			playlistController.GetUserPlaylists)
		playlist.DELETE("/delete-tracks",
			middlewares.SpotifyAuthMiddlewareRefactored(),
			playlistController.DeleteAllPlaylistTracks)
		playlist.DELETE("/delete-tracks-and-library",
			middlewares.SpotifyAuthMiddlewareRefactored(),
			playlistController.DeleteAllPlaylistAndUserTracks)
	}
}
