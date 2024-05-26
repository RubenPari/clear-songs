package routes

import (
	"github.com/RubenPari/clear-songs/src/controllers"
	"github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	// ####### AUTHENTICATION #######
	auth := server.Group("/auth")
	{
		auth.GET("/login", controllers.Login)
		auth.GET("/callback", controllers.Callback)
		auth.GET("/logout", controllers.Logout)
	}

	// ####### TRACK #######
	track := server.Group("/track")
	{
		track.GET("/summary",
			middlewares.CheckAuth(),
			controllers.GetTrackSummary)
		track.DELETE("/by-artist/:id_artist",
			middlewares.CheckAuth(),
			controllers.DeleteTrackByArtist)
		track.DELETE("/by-range",
			middlewares.CheckAuth(),
			controllers.DeleteTrackByRange)
	}

	// ####### ALBUMS #######
	album := server.Group("/album")
	{
		album.GET("/all",
			middlewares.CheckAuth(),
			controllers.GetAll)
		album.GET("/by-artist/:id_artist",
			middlewares.CheckAuth(),
			controllers.GetAlbumByArtist)
		album.PUT("/convert-to-songs",
			middlewares.CheckAuth(),
			controllers.ConvertAlbumToSongs)
	}

	// ####### PLAYLIST #######
	playlist := server.Group("/playlist")
	{
		playlist.DELETE("/delete-tracks",
			middlewares.CheckAuth(),
			controllers.DeleteAllPlaylistTracks)
	}
}
