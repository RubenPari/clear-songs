package routes

import (
	"github.com/RubenPari/clear-songs/src/controllers/albumController"
	"github.com/RubenPari/clear-songs/src/controllers/authController"
	"github.com/RubenPari/clear-songs/src/controllers/playlistController"
	"github.com/RubenPari/clear-songs/src/controllers/trackController"
	"github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	// 404 Not Found Route
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
		track.DELETE("/by-artist/:id_artist",
			middlewares.CheckAuth(),
			trackController.DeleteTrackByArtist)
		track.DELETE("/by-range",
			middlewares.CheckAuth(),
			trackController.DeleteTrackByRange)
	}

	// ####### ALBUMS #######
	server.POST("/album/convert-to-songs",
		middlewares.CheckAuth(),
		albumController.ConvertAlbumToSongs)

	// ####### PLAYLIST #######
	playlist := server.Group("/playlist")
	{
		playlist.DELETE("/delete-tracks",
			middlewares.CheckAuth(),
			playlistController.DeleteAllPlaylistTracks)
		playlist.DELETE("/delete-tracks-and-library",
			middlewares.CheckAuth(),
			playlistController.DeleteAllPlaylistAndUserTracks)
	}
}
