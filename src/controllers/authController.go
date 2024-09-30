package controllers

import (
	"context"
	"log"

	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// Login gestisce la richiesta di login dell'utente.
//
// Effettua la seguente logica:
// 1. Ottiene la configurazione OAuth2.
// 2. Crea l'URL di login per Spotify.
// 3. Reindirizza l'utente all'URL di login.
func Login(c *gin.Context) {
	configAuth := utils.GetOAuth2Config()

	url := configAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Default().Printf("Redirecting to %s", url)

	c.Redirect(302, url)
}

// Callback gestisce la callback ricevuta da Spotify in seguito alla richiesta
// di login dell'utente.
//
// Effettua la seguente logica:
//  1. Ottiene il codice di autorizzazione ricevuto nella query string.
//  2. Effettua l'exchange del codice con il token di accesso.
//  3. Crea un client Spotify utilizzando il token di accesso.
//  4. Salva il client Spotify nella sessione.
//  5. Esegue una chiamata di test per verificare l'autenticazione.
//  6. Restituisce un JSON di successo se l'autenticazione va a buon fine,
//     altrimenti restituisce un JSON di errore.
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

// Logout cancella l'autenticazione dell'utente cancellando il client Spotify
// dalla sessione.
//
// Effettua la seguente logica:
//  1. Cancella il client Spotify dalla sessione.
//  2. Restituisce un JSON di successo se l'autenticazione va a buon fine,
//     altrimenti restituisce un JSON di errore.
func Logout(c *gin.Context) {
	// delete spotify client from session
	utils.SpotifyClient = spotifyAPI.Client{}

	log.Default().Println("Called logout, deleted client from session")

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User logged out",
	})
}

// IsAuth verifica se l'utente è autenticato.
//
// Effettua la seguente logica:
//  1. Verifica se il client Spotify è settato.
//  2. Se il client Spotify non è settato, restituisce un JSON di errore con
//     lo stato 401 e un messaggio di errore.
//  3. Se il client Spotify è settato, restituisce un JSON di successo con
//     lo stato 200 e un messaggio di successo.
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
