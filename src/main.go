package main

import (
	"github.com/RubenPari/clear-songs/src/cacheManager"
	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize server
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// set routes
	routes.SetUpRoutes(server)

	// init cache
	cacheManager.Init()

	// connect to database
	if errConnectDb := database.Init(); errConnectDb != nil {
		panic("Error connecting to database")
	}

	// start server
	if server.Run("0.0.0.0:8080") != nil {
		panic("Error starting server")
	}
}
