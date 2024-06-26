package main

import (
	"github.com/RubenPari/clear-songs/src/cacheManager"
	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	// initialize server
	gin.SetMode(gin.DebugMode)
	server := gin.Default()

	// set routes
	routes.SetUpRoutes(server)

	// load environment variables
	utils.LoadEnvVariables()

	// init cache
	cacheManager.Init()

	// connect to database
	if errConnectDb := database.Init(); errConnectDb != nil {
		panic("Error connecting to database")
	}

	// start server
	if server.Run(":"+os.Getenv("PORT")) != nil {
		panic("Error starting server")
	}
}
