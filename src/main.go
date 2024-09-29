package main

import (
	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/database"
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

	// load environment variables
	utils.LoadEnvVariables()

	// init cache
	cacheManager.Init()

	// connect to database
	if errConnectDb := database.Init(); errConnectDb != nil {
		panic("Error connecting to database")
	}

	// start server
	if server.Run(":8080") != nil {
		panic("Error starting server")
	}
}
