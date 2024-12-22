package main

import (
	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/docs"
	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Clear Songs API
// @version 1.0
// @description API for managing Spotify playlists and tracks
// @BasePath /
func main() {
	gin.SetMode(gin.DebugMode)
	server := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
