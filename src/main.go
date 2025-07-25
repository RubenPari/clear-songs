package main

import (
	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	server := gin.Default()


	routes.SetUpRoutes(server)

	utils.LoadEnvVariables()

	cacheManager.Init()

	if errConnectDb := database.Init(); errConnectDb != nil {
		panic("Error connecting to database")
	}

	if server.Run(":3000") != nil {
		panic("Error starting server")
	}
}
