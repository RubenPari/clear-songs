package main

import (
	"os"

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

	// connect to database
	if errConnectDb := database.Init(); errConnectDb != nil {
		panic("Error connecting to database")
	}

	// start server
	if server.Run(":"+os.Getenv("PORT")) != nil {
		panic("Error starting server")
	}
}
