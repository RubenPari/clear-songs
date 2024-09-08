package controllers

import (
	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
	"log"
)

func Logout(c *gin.Context) {
	// set spotify client to nil
	utils.SpotifyClient = spotifyAPI.Client{}

	// flash session
	cacheManager.Delete()

	log.Default().Println("Called logout, deleted client from session")

	c.JSON(200, gin.H{
		"message": "User logged out",
	})
}
