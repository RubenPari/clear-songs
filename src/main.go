package main

import (
	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize server
	gin.SetMode(gin.DebugMode)
	server := gin.Default()

	// set routes
	routes.SetUpRoutes(server)

	// start server
	if server.Run(":"+utils.Port) != nil {
		panic("Error starting server")
	}
}
