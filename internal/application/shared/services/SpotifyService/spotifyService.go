/**
 * Spotify Service Package
 * 
 * This package provides a wrapper around the Spotify Web API client library.
 * It manages Spotify OAuth authentication and provides a configured client
 * instance for making authenticated API calls to Spotify.
 * 
 * The service handles:
 * - OAuth 2.0 authentication configuration
 * - Access token management
 * - Spotify client initialization
 * - API call authorization
 * 
 * The service is designed as a singleton that can be shared across the
 * application, with tokens being updated per-user session.
 * 
 * @package SpotifyService
 * @author Clear Songs Development Team
 */
package SpotifyService

import (
	"github.com/RubenPari/clear-songs/internal/domain/shared/constants"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

/**
 * SpotifyService struct
 * 
 * Encapsulates Spotify API client configuration and authentication state.
 * This struct holds all necessary information to make authenticated calls
 * to the Spotify Web API.
 * 
 * Fields:
 * - clientID: Spotify application client ID
 * - clientSecret: Spotify application client secret
 * - redirectURI: OAuth callback URL registered with Spotify
 * - authenticator: Spotify OAuth authenticator instance
 * - token: Current OAuth access token (per user session)
 * - client: Authenticated Spotify API client
 */
type SpotifyService struct {
	clientID      string
	clientSecret  string
	redirectURI   string
	authenticator spotify.Authenticator
	token         *oauth2.Token
	client        *spotify.Client
}

/**
 * NewSpotifyService creates a new Spotify service instance
 * 
 * Initializes a SpotifyService with OAuth configuration. The service
 * is configured with the required OAuth scopes defined in constants
 * and is ready to handle authentication flows.
 * 
 * @param clientID - Spotify application client ID from developer dashboard
 * @param clientSecret - Spotify application client secret
 * @param redirectURI - OAuth callback URL (must match Spotify app settings)
 * @returns *SpotifyService - Configured service instance ready for use
 */
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
