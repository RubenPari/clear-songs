package controllers

import (
	"context"
	"github.com/RubenPari/clear-songs/src/cacheManager"
	"log"

	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var configAuth = utils.GetOAuth2Config()

func Login(c *gin.Context) {
	// create url for spotify login
	url := configAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Default().Println("Called login, redirecting to: " + url)

	// redirect to url
	c.Redirect(302, url)
}

func Callback(c *gin.Context) {
	code := c.Query("code")

	// get token from code
	token, errToken := configAuth.Exchange(context.Background(), code)

	if errToken != nil {
		c.JSON(500, gin.H{
			"message": "Error authenticating user",
		})
	}

	// create client
	client := configAuth.Client(context.Background(), token)
	spotify := spotifyAPI.NewClient(client)

	log.Default().Println("Called callback, created client")

	// save spotify client in session
	utils.SpotifyClient = &spotify

	// get user for testing
	user, errUser := spotify.CurrentUser()

	if errUser != nil {
		c.JSON(500, gin.H{
			"message": "Error authenticating user",
		})
	}

	log.Default().Println("Called callback from user", user.User.DisplayName)

	c.JSON(200, gin.H{
		"message": "User authenticated",
	})
}

func Logout(c *gin.Context) {
	// set spotify client to nil
	utils.SpotifyClient = nil

	// flash session
	cacheManager.Delete()

	log.Default().Println("Called logout, deleted client from session")

	c.JSON(200, gin.H{
		"message": "User logged out",
	})
}
