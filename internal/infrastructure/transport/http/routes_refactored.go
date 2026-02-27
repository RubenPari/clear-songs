package http

import (
	"github.com/RubenPari/clear-songs/internal/infrastructure/di"
	"github.com/RubenPari/clear-songs/internal/infrastructure/transport/http/handlers"
	"github.com/RubenPari/clear-songs/internal/infrastructure/transport/http/middleware"
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
	server.Use(middleware.SessionMiddlewareRefactored(
		container.SpotifyRepo,
		container.CacheRepo,
	))
	server.Use(middleware.CacheInvalidationMiddleware())

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
	 */
	trackController := handlers.NewTrackControllerComplete(
		container.GetTrackSummaryUseCase,
		container.DeleteTracksByArtistUC,
		container.DeleteTracksByRangeUC,
		container.GetTracksByArtistUC,
		container.DeleteTrackUC,
	)

	track := server.Group("/track")
	{
		track.GET("/summary",
			middleware.SpotifyAuthMiddlewareRefactored(),
			trackController.GetTrackSummary)
		track.GET("/by-artist/:id_artist",
			middleware.SpotifyAuthMiddlewareRefactored(),
			trackController.GetTracksByArtist)
		track.DELETE("/by-artist/:id_artist",
			middleware.SpotifyAuthMiddlewareRefactored(),
			trackController.DeleteTrackByArtist)
		track.DELETE("/:id_track",
			middleware.SpotifyAuthMiddlewareRefactored(),
			trackController.DeleteTrack)
		track.DELETE("/by-range",
			middleware.SpotifyAuthMiddlewareRefactored(),
			trackController.DeleteTrackByRange)
	}

	/**
	 * Authentication Routes Group
	 */
	authController := handlers.NewAuthControllerRefactored(
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
	 */
	playlistController := handlers.NewPlaylistControllerRefactored(
		container.GetUserPlaylistsUC,
		container.DeletePlaylistTracksUC,
		container.DeletePlaylistAndLibraryUC,
	)

	playlist := server.Group("/playlist")
	{
		playlist.GET("/list",
			middleware.SpotifyAuthMiddlewareRefactored(),
			playlistController.GetUserPlaylists)
		playlist.DELETE("/delete-tracks",
			middleware.SpotifyAuthMiddlewareRefactored(),
			playlistController.DeleteAllPlaylistTracks)
		playlist.DELETE("/delete-tracks-and-library",
			middleware.SpotifyAuthMiddlewareRefactored(),
			playlistController.DeleteAllPlaylistAndUserTracks)
	}
}
