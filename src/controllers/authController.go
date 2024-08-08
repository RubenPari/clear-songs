package controllers

import (
	"context"
	"log"

	"github.com/RubenPari/clear-songs/src/cacheManager"
	"github.com/RubenPari/clear-songs/src/models"

	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var configAuth *oauth2.Config = nil

func LoginApi(c *gin.Context) {
	configAuth = utils.GetOAuth2Config()

	// create url for spotify login
	url := configAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Default().Println("Called login, redirecting to: " + url)

	// redirect to url
	c.Redirect(302, url)
}

func LoginFront(c *gin.Context) {
	var accessTokenRequest models.AccessTokenRequest

	// parse access token from body request already got from front-end
	errParsingBody := c.ShouldBindJSON(&accessTokenRequest)

	if errParsingBody != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// set token
	token := &oauth2.Token{
		AccessToken: accessTokenRequest.AccessToken,
		TokenType:   "Bearer",
	}

	// create client
	client := configAuth.Client(context.Background(), token)
	spotify := spotifyAPI.NewClient(client)

	// save spotify client in session
	utils.SpotifyClient = &spotify

	// get user for testing
	user, errUser := utils.SpotifyClient.CurrentUser()

	if errUser != nil {
		c.JSON(500, gin.H{
			"message": "Error authenticating user",
		})
	}

	log.Default().Println("Called login from user", user.User.DisplayName)

	c.JSON(200, gin.H{
		"message": "User authenticated",
	})

}

func Callback(c *gin.Context) {
	code := c.Query("code")

	// get token from code
	token, errToken := configAuth.Exchange(context.Background(), code)

	if errToken != nil {
		c.JSON(500, gin.H{
			"message": "Error authenticating user",
			"error":   errToken.Error(),
		})
	}

	// create client
	client := configAuth.Client(context.Background(), token)
	spotify := spotifyAPI.NewClient(client)

	log.Default().Println("Called callback, created client")

	// save spotify client in session
	utils.SpotifyClient = &spotify

	log.Default().Println("Called callback, saved client in session")

	// get user for testing
	user, errUser := utils.SpotifyClient.CurrentUser()

	if errUser != nil {
		c.JSON(500, gin.H{
			"message": "Error getting user info",
			"error":   errUser.Error(),
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
