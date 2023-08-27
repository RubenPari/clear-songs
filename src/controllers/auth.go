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

	// generate randomm token for protected ednpoints
	tokenHeader := utils.RandomString(20)

	// save generated token in session
	utils.TokenHeader = tokenHeader

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
	utils.SpotifyClient = &spotify

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
		"token":   tokenHeader,
	})
}

func Logout(c *gin.Context) {
	// delete spotify client from session
	utils.SpotifyClient = nil
	utils.TokenHeader = ""

	log.Default().Println("Called logout, deleted client from session")

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User logged out",
	})
}
