package middlewares

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/RubenPari/clear-songs/src/models"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from microservice dedicated to log in
		response, errResponse := http.Get(os.Getenv("MICROSERVICE_AUTH_LOGIN_SPOTIFY") + "/token")

		if errResponse != nil {
			log.Default().Printf("Errore durante la richiesta del token: %v", errResponse)

			c.AbortWithStatusJSON(500, gin.H{
				"message": "Errore durante la richiesta del token",
			})

			return
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(response.Body)

		// Read body of response
		tokenBytes, errTokenBytes := io.ReadAll(response.Body)

		if errTokenBytes != nil {
			log.Default().Printf("Errore durante la lettura del corpo della risposta: %v", errTokenBytes)

			c.AbortWithStatusJSON(500, gin.H{
				"message": "Errore durante la lettura del corpo della risposta",
			})

			return
		}

		// Unmarshal il token
		var accessToken models.AccessTokenResponse

		errUnmarshalToken := json.Unmarshal(tokenBytes, &accessToken)

		if errUnmarshalToken != nil {
			log.Default().Printf("Errore durante la serializzazione del token in json: %v", errUnmarshalToken)

			c.AbortWithStatusJSON(500, gin.H{
				"message": "Errore durante la serializzazione del token in json",
			})

			return
		}

		// Create a new Spotify client
		spotifyClient := spotify.Authenticator{}.NewClient(&oauth2.Token{
			AccessToken: accessToken.AccessToken,
		})

		// Check if works with the token
		_, errCurrentUser := spotifyClient.CurrentUser()

		if errCurrentUser != nil {
			log.Default().Printf("Errore durante la verifica del token: %v", errCurrentUser)

			c.AbortWithStatusJSON(500, gin.H{
				"message": "Errore durante la verifica del token",
			})

			return
		}

		// Save the spotify client into the context
		utils.SpotifyClient = spotifyClient

		c.Next()
	}
}
