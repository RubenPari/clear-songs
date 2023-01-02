package auth

import (
	"os"
	"strings"

	spotifyAPI "github.com/zmb3/spotify/v2"

	"golang.org/x/oauth2"
	spotifyOAuth "golang.org/x/oauth2/spotify"
)

var (
	SpotifyClient *spotifyAPI.Client = nil
	Token         *oauth2.Token      = nil
)

func GenerateStateString() string {
	// TODO: generate random string
	return "state"
}

func GetOAuthConfig() (*oauth2.Config, string) {
	scopes := os.Getenv("SCOPES")
	scopesArray := strings.Split(scopes, ",")

	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       scopesArray,
		Endpoint:     spotifyOAuth.Endpoint,
	}, GenerateStateString()
}
