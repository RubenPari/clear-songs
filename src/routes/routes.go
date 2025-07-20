package routes

import (
	"github.com/RubenPari/clear-songs/src/controllers/authController"
	"github.com/RubenPari/clear-songs/src/controllers/playlistController"
	"github.com/RubenPari/clear-songs/src/controllers/trackController"
	"github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	server.Use(middlewares.SessionMiddleware())
	server.Use(middlewares.CacheInvalidationMiddleware())

	// ####### NOT FOUND ROUTE #######
	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "not found path",
		})
	})

	// ####### AUTHENTICATION #######
	auth := server.Group("/auth")
	{
		auth.GET("/login", authController.Login)
		auth.GET("/callback", authController.Callback)
		auth.GET("/logout", authController.Logout)
		auth.GET("/is-auth", authController.IsAuth)
	}

	// ####### TRACK #######
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

	// ####### PLAYLIST #######
	playlist := server.Group("/playlist")
	{
		playlist.DELETE("/delete-tracks",
			middlewares.SpotifyAuthMiddleware(),
			playlistController.DeleteAllPlaylistTracks)
		playlist.DELETE("/delete-tracks-and-library",
			middlewares.SpotifyAuthMiddleware(),
			playlistController.DeleteAllPlaylistAndUserTracks)
	}
}
