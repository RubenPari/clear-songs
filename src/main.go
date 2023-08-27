package main

import (
	"os"

	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize server
	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// set routes
	routes.SetUpRoutes(server)

	// set CORS
	utils.CorsConfig(server)

	// start server
	if server.Run(":"+os.Getenv("PORT")) != nil {
		panic("Error starting server")
	}
}
