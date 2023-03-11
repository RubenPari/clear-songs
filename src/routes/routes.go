package routes

import (
	authContr "github.com/RubenPari/clear-songs/src/controllers"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	// ####### AUTHENTICATION #######
	server.GET("/auth/login", authContr.Login)
	server.GET("/auth/callback", authContr.Callback)

	// ####### TRACK #######
	server.GET("/track/summary", authContr.GetTrackSummary)
	server.DELETE("/track/by-aritst/:id_artist", authContr.DeleteTrackByArtist)
	server.DELETE("/track/by-genre", authContr.DeleteTrackByGenre)
}
