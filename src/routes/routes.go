package routes

import (
	authContr "github.com/RubenPari/clear-songs/src/controllers"
	"github.com/RubenPari/clear-songs/src/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	// ####### AUTHENTICATION #######
	server.GET("/auth/login", authContr.Login)
	server.GET("/auth/callback", authContr.Callback)
	server.GET("/auth/logout", authContr.Logout)

	// ####### TRACK #######
	server.GET("/track/summary",
		middlewares.CheckAuth(),
		authContr.GetTrackSummary)
	server.DELETE("/track/by-artist/:id_artist",
		middlewares.CheckAuth(),
		authContr.DeleteTrackByArtist)
	server.DELETE("/track/by-genre",
		middlewares.CheckAuth(),
		authContr.DeleteTrackByGenre)
	server.DELETE("/track/by-range",
		middlewares.CheckAuth(),
		authContr.DeleteTrackByRange)
}
