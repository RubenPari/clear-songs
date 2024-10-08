package controllers

import (
	"context"
	"log"

	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func Login(c *gin.Context) {
	// get oauth2 config
	configAuth := utils.GetOAuth2Config()

	// create url for spotify login
	url := configAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Default().Println("Called login, redirecting to: " + url)

	// redirect to url
	c.Redirect(302, url)
}

func Callback(c *gin.Context) {
	// get code from query parameters
	code := c.Query("code")

	// get oauth2 config
	configAuth := utils.GetOAuth2Config()

	// get token from code
	token, errToken := configAuth.Exchange(context.Background(), code)

	if errToken != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error authenticating user",
		})
	}

	// create client
	client := configAuth.Client(context.Background(), token)
	spotify := spotifyAPI.NewClient(client)

	log.Default().Println("Called callback, created client")

	// save spotify client in session
	utils.SpotifyClient = spotify

	// get user for testing
	_, errUser := spotify.CurrentUser()

	if errUser != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error authenticating user",
		})
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User authenticated",
	})
}

func Logout(c *gin.Context) {
	// delete spotify client from session
	utils.SpotifyClient = spotifyAPI.Client{}

	log.Default().Println("Called logout, deleted client from session")

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User logged out",
	})
}

func IsAuth(c *gin.Context) {
	// check if spotify client is set
	if _, err := utils.SpotifyClient.CurrentUser(); err != nil {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "User authenticated",
		})
	}
}
