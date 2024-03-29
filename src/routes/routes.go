package routes

import (
	albumContr "github.com/RubenPari/clear-songs/src/controllers"
	authContr "github.com/RubenPari/clear-songs/src/controllers"
	utilsContr "github.com/RubenPari/clear-songs/src/controllers"
	"github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	// ####### AUTHENTICATION #######
	auth := server.Group("/auth")
	{
		auth.GET("/login", authContr.Login)
		auth.GET("/callback", authContr.Callback)
		auth.GET("/logout", authContr.Logout)
	}

	// ####### TRACK #######
	track := server.Group("/track")
	{
		track.GET("/summary",
			middlewares.CheckAuth(),
			authContr.GetTrackSummary)
		track.DELETE("/by-artist/:id_artist",
			middlewares.CheckAuth(),
			authContr.DeleteTrackByArtist)
		track.DELETE("/by-range",
			middlewares.CheckAuth(),
			authContr.DeleteTrackByRange)
		track.DELETE("/by-file",
			middlewares.CheckAuth(),
			authContr.DeleteTrackByFile)
	}

	// ####### ALBUMS #######
	album := server.Group("/album")
	{
		album.GET("/all",
			middlewares.CheckAuth(),
			albumContr.GetAll)
		album.GET("/by-artist/:id_artist",
			middlewares.CheckAuth(),
			albumContr.GetAlbumByArtist)
		album.DELETE("/by-artist/:id_artist",
			middlewares.CheckAuth(),
			albumContr.DeleteAlbumByArtist)
		album.PUT("/convert-to-songs",
			middlewares.CheckAuth(),
			albumContr.ConvertAlbumToSongs)
	}

	// ####### UTILS #######
	utils := server.Group("/utils")
	{
		utils.GET("/name-by-id/:id",
			middlewares.CheckAuth(),
			utilsContr.GetNameByID)
		utils.GET("/id-by-name",
			middlewares.CheckAuth(),
			utilsContr.GetIDByName)
	}
}
