package main

import (
	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	// initialize server
	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// load .env file
	utils.LoadEnv(1)

	// set routes
	routes.SetUpRoutes(server)

	// start server
	if server.Run(":"+os.Getenv("PORT")) != nil {
		panic("Error starting server")
	}
}
