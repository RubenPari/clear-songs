package SpotifyService

import (
	"github.com/RubenPari/clear-songs/src/constants"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type SpotifyService struct {
	clientID      string
	clientSecret  string
	redirectURI   string
	authenticator spotify.Authenticator
	token         *oauth2.Token
	client        *spotify.Client
}

func NewSpotifyService(clientID, clientSecret, redirectURI string) *SpotifyService {
	auth := spotify.NewAuthenticator(redirectURI, constants.Scopes...)
	auth.SetAuthInfo(clientID, clientSecret)

	return &SpotifyService{
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURI:   redirectURI,
		authenticator: auth,
	}
}

func (s *SpotifyService) SetAccessToken(token *oauth2.Token) {
	s.token = token
	client := s.authenticator.NewClient(token)
	s.client = &client
}

func (s *SpotifyService) GetSpotifyClient() *spotify.Client {
	return s.client
}
