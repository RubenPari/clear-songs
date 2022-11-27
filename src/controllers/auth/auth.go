package auth

import (
	"context"
	"log"
	"net/http"

	spotifyAPI "github.com/zmb3/spotify/v2"
	spotifyAUTH "github.com/zmb3/spotify/v2/auth"

	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/gofiber/fiber/v2"
)

var (
	oauthConf, stateGlobal = authMO.GetOAuthConfig()
)

func Login(c *fiber.Ctx) error {
	url := oauthConf.AuthCodeURL(stateGlobal)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func Callback(c *fiber.Ctx) error {
	state := c.Query("state")

	if state != stateGlobal {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", stateGlobal, state)
		_ = c.SendStatus(http.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"status": "error",
			"error":  "invalid oauth state",
		})
	}

	code := c.Query("code")
	ctx := context.Background()

	token, err := oauthConf.Exchange(ctx, code)

	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		_ = c.SendStatus(http.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"status": "error",
			"error":  "Code exchange failed",
		})
	}

	// create http client
	httpClient := spotifyAUTH.New().Client(ctx, token)

	// export token
	authMO.Token = token

	// create spotify http client
	spotifyClient := spotifyAPI.New(httpClient)

	// export spotify client
	authMO.SpotifyClient = spotifyClient

	resp, err := spotifyClient.CurrentUser(ctx)

	if err != nil {
		_ = c.SendStatus(http.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "Failed authenticating",
		})
	}

	_ = c.SendStatus(http.StatusOK)
	return c.JSON(fiber.Map{
		"status":       "success",
		"message":      "login successful",
		"current_user": resp,
	})
}
