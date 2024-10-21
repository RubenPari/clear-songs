package controllers

import (
	"context"
	"log"

	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// Login redirects to Spotify's authentication address.
//
// The function uses the value of oauth2.AccessTypeOffline to get
// an offline access token that can be used to make
// API calls in the future without having to prompt the user to authenticate
// again.
func Login(c *gin.Context) {
	configAuth := utils.GetOAuth2Config()

	url := configAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Default().Printf("Redirecting to %s", url)

	c.Redirect(302, url)
}

// Callback handles the callback received from Spotify after the user
// login request.
//
// Performs the following logic:
// 1. Gets the authorization code received in the query string.
// 2. Exchanges the code with the access token.
// 3. Creates a Spotify client using the access token.
// 4. Saves the Spotify client to the session.
// 5. Makes a test call to verify authentication.
// 6. Returns a success JSON if authentication is successful,
// otherwise returns an error JSON.
func Callback(c *gin.Context) {
	// get code from query parameters
	code := c.Query("code")

	configAuth := utils.GetOAuth2Config()

	token, errToken := configAuth.Exchange(context.Background(), code)

	if errToken != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error authenticating user",
		})
	}

	// create spotify client
	client := configAuth.Client(context.Background(), token)
	spotify := spotifyAPI.NewClient(client)

	log.Default().Println("Called callback, created spotify wrapper")

	// save spotify client in session
	utils.SpotifyClient = spotify

	// get user info for testing
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

// Logout clears the user authentication by clearing the Spotify client
// from the session.
//
// Performs the following logic:
// 1. Clears the Spotify client from the session.
// 2. Returns a success JSON if authentication is successful,
// otherwise returns an error JSON.
func Logout(c *gin.Context) {
	// delete spotify client from session
	utils.SpotifyClient = spotifyAPI.Client{}

	log.Default().Println("Called logout, deleted client from session")

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User logged out",
	})
}

// IsAuth checks if the user is authenticated.
//
// Performs the following logic:
// 1. Checks if Spotify client is set.
// 2. If Spotify client is not set, returns an error JSON with
// status 401 and an error message.
// 3. If Spotify client is set, returns a success JSON with
// status 200 and a success message.
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
