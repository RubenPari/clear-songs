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

	// load .env file
	utils.LoadEnv()

	// set routes
	routes.SetUpRoutes(server)

	// start server
	if server.Run(":"+os.Getenv("PORT")) != nil {
		panic("Error starting server")
	}
}
