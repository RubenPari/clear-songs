package routes

import (
	"github.com/RubenPari/clear-songs/src/controllers"
	"github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	server.GET("/", middlewares.NotFound())

	api := server.Group("/api")
	apiV1 := api.Group("/v1")

	// TODO: simplify auth
	// ####### AUTHENTICATION #######
	auth := apiV1.Group("/auth")
	{
		auth.GET("/login-api", controllers.LoginApi)
		auth.GET("/login-front", controllers.LoginFront)
		auth.GET("/callback", controllers.Callback)
		auth.GET("/logout", controllers.Logout)
	}

	// TODO: remove summary
	// ####### TRACK #######
	track := apiV1.Group("/track")
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

	// TODO: microservices apart from
	// ####### ALBUMS #######
	album := apiV1.Group("/album")
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
	playlist := apiV1.Group("/playlist")
	{
		playlist.DELETE("/delete-tracks",
			middlewares.CheckAuth(),
			controllers.DeleteAllPlaylistTracks)
		playlist.DELETE("/delete-tracks-and-library",
			middlewares.CheckAuth(),
			controllers.DeleteAllPlaylistAndUserTracks)
	}
}
