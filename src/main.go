package main

import (
	"os"

	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize server
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// set routes
	routes.SetUpRoutes(server)

	// start server
	if server.Run(":"+os.Getenv("PORT")) != nil {
		panic("Error starting server")
	}
}
