package main

import (
	"github.com/RubenPari/clear-songs/src/lib/utils"
	"os"

	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize server
	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// load .env file
	utils.LoadEnv(0)

	// set routes
	routes.SetUpRoutes(server)

	// start server
	if server.Run(":"+os.Getenv("PORT")) != nil {
		panic("Error starting server")
	}
}
