package routes

import (
	controllers2 "github.com/RubenPari/clear-songs/src/controllers"
	middlewares2 "github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	server.GET("/", middlewares2.NotFound())

	// ####### AUTHENTICATION #######
	auth := server.Group("/auth")
	{
		auth.GET("/login", controllers2.Login)
		auth.GET("/callback", controllers2.Callback)
		auth.GET("/logout", controllers2.Logout)
	}

	// ####### TRACK #######
	track := server.Group("/track")
	{
		track.GET("/summary",
			middlewares2.CheckAuth(),
			controllers2.GetTrackSummary)
		track.DELETE("/by-artist/:id_artist",
			middlewares2.CheckAuth(),
			controllers2.DeleteTrackByArtist)
		track.DELETE("/by-range",
			middlewares2.CheckAuth(),
			controllers2.DeleteTrackByRange)
	}

	// ####### ALBUMS #######
	album := server.Group("/album")
	{
		album.GET("/all",
			middlewares2.CheckAuth(),
			controllers2.GetAll)
		album.GET("/by-artist/:id_artist",
			middlewares2.CheckAuth(),
			controllers2.GetAlbumByArtist)
		album.PUT("/convert-to-songs",
			middlewares2.CheckAuth(),
			controllers2.ConvertAlbumToSongs)
	}

	// ####### PLAYLIST #######
	playlist := server.Group("/playlist")
	{
		playlist.DELETE("/delete-tracks",
			middlewares2.CheckAuth(),
			controllers2.DeleteAllPlaylistTracks)
		playlist.DELETE("/delete-tracks-and-library",
			middlewares2.CheckAuth(),
			controllers2.DeleteAllPlaylistAndUserTracks)
	}
}
