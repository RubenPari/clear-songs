package routes

import (
	authContr "github.com/RubenPari/clear-songs/src/controllers"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server *gin.Engine) {
	// ####### AUTHENTICATION #######
	server.GET("/auth/login", authContr.Login)
	server.GET("/auth/callback", authContr.Callback)
}
