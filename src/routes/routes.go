package routes

import (
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
		track.DELETE("/by-genre",
			middlewares.CheckAuth(),
			authContr.DeleteTrackByGenre)
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
			authContr.GetAll)
		album.GET("/by-artist/:id_artist",
			middlewares.CheckAuth(),
			authContr.GetAlbumByArtist)
		album.DELETE("/by-artist/:id_artist",
			middlewares.CheckAuth(),
			authContr.DeleteAlbumByArtist)
		album.PUT("/convert-to-songs",
			middlewares.CheckAuth(),
			authContr.ConvertAlbumToSongs)
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
