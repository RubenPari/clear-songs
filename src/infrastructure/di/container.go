package di

import (
	"errors"
	"log"
	"os"

	"github.com/RubenPari/clear-songs/src/application/usecases/auth"
	"github.com/RubenPari/clear-songs/src/application/usecases/playlist"
	"github.com/RubenPari/clear-songs/src/application/usecases/track"
	"github.com/RubenPari/clear-songs/src/constants"
	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/domain/interfaces"
	"github.com/RubenPari/clear-songs/src/infrastructure/cache"
	dbRepo "github.com/RubenPari/clear-songs/src/infrastructure/database"
	"github.com/RubenPari/clear-songs/src/infrastructure/spotify"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// Container holds all application dependencies
type Container struct {
	// Repositories (as interfaces)
	SpotifyRepo  interfaces.SpotifyRepository
	CacheRepo    interfaces.CacheRepository
	DatabaseRepo interfaces.DatabaseRepository
	
	// OAuth Config
	OAuthConfig *oauth2.Config
	
	// Auth Use Cases
	LoginUC    *auth.LoginUseCase
	CallbackUC *auth.CallbackUseCase
	LogoutUC   *auth.LogoutUseCase
	IsAuthUC   *auth.IsAuthUseCase
	
	// Track Use Cases
	GetTrackSummaryUseCase    *track.GetTrackSummaryUseCase
	DeleteTracksByArtistUC    *track.DeleteTracksByArtistUseCase
	DeleteTracksByRangeUC     *track.DeleteTracksByRangeUseCase
	GetTracksByArtistUC       *track.GetTracksByArtistUseCase
	
	// Playlist Use Cases
	GetUserPlaylistsUC         *playlist.GetUserPlaylistsUseCase
	DeletePlaylistTracksUC     *playlist.DeletePlaylistTracksUseCase
	DeletePlaylistAndLibraryUC *playlist.DeletePlaylistAndLibraryTracksUseCase
}

// NewContainer creates and initializes a new dependency injection container
func NewContainer() (*Container, error) {
	// Initialize Spotify repository
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URL")
	if redirectURI == "" {
		redirectURI = os.Getenv("REDIRECT_URI")
	}
	
	if clientID == "" || clientSecret == "" || redirectURI == "" {
		log.Fatal("Missing required environment variables: CLIENT_ID, CLIENT_SECRET, REDIRECT_URL")
	}
	
	spotifyRepo := spotify.NewSpotifyRepository(clientID, clientSecret, redirectURI, constants.Scopes)
	
	// Initialize OAuth config
	oauthConfig, err := GetOAuth2Config()
	if err != nil {
		return nil, err
	}
	
	// Initialize cache repository (may fail if Redis is not available)
	var cacheRepo interfaces.CacheRepository
	redisCache, err := cache.NewRedisCacheRepository()
	if err != nil {
		log.Printf("WARNING: Cache repository initialization failed: %v", err)
		log.Println("WARNING: Application will continue without Redis caching")
		// Use no-op cache implementation
		cacheRepo = cache.NewNoOpCacheRepository()
	} else {
		cacheRepo = redisCache
	}
	
	// Initialize auth use cases
	loginUC := auth.NewLoginUseCase(oauthConfig)
	callbackUC := auth.NewCallbackUseCase(oauthConfig, spotifyRepo, cacheRepo)
	logoutUC := auth.NewLogoutUseCase(spotifyRepo, cacheRepo)
	isAuthUC := auth.NewIsAuthUseCase(spotifyRepo)
	
	// Initialize database repository (may be nil if database not available)
	databaseRepo := dbRepo.NewPostgresRepository(database.Db)
	
	// Initialize track use cases
	getTrackSummaryUseCase := track.NewGetTrackSummaryUseCase(spotifyRepo, cacheRepo)
	deleteTracksByArtistUC := track.NewDeleteTracksByArtistUseCase(spotifyRepo, cacheRepo)
	getTracksByArtistUC := track.NewGetTracksByArtistUseCase(spotifyRepo, cacheRepo)
	deleteTracksByRangeUC := track.NewDeleteTracksByRangeUseCase(
		spotifyRepo,
		cacheRepo,
		getTrackSummaryUseCase,
		deleteTracksByArtistUC,
	)
	
	// Initialize playlist use cases
	getUserPlaylistsUC := playlist.NewGetUserPlaylistsUseCase(spotifyRepo, cacheRepo)
	deletePlaylistTracksUC := playlist.NewDeletePlaylistTracksUseCase(spotifyRepo, cacheRepo)
	deletePlaylistAndLibraryUC := playlist.NewDeletePlaylistAndLibraryTracksUseCase(
		spotifyRepo,
		cacheRepo,
		databaseRepo,
		deletePlaylistTracksUC,
	)
	
	container := &Container{
		SpotifyRepo:               spotifyRepo,
		CacheRepo:                 cacheRepo,
		DatabaseRepo:              databaseRepo,
		OAuthConfig:               oauthConfig,
		LoginUC:                   loginUC,
		CallbackUC:                callbackUC,
		LogoutUC:                  logoutUC,
		IsAuthUC:                  isAuthUC,
		GetTrackSummaryUseCase:    getTrackSummaryUseCase,
		DeleteTracksByArtistUC:    deleteTracksByArtistUC,
		DeleteTracksByRangeUC:    deleteTracksByRangeUC,
		GetTracksByArtistUC:       getTracksByArtistUC,
		GetUserPlaylistsUC:        getUserPlaylistsUC,
		DeletePlaylistTracksUC:    deletePlaylistTracksUC,
		DeletePlaylistAndLibraryUC: deletePlaylistAndLibraryUC,
	}
	
	return container, nil
}

// GetOAuth2Config returns OAuth2 configuration from environment variables
func GetOAuth2Config() (*oauth2.Config, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URL")
	if redirectURI == "" {
		redirectURI = os.Getenv("REDIRECT_URI")
	}
	
	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, errors.New("missing required environment variables: CLIENT_ID, CLIENT_SECRET, REDIRECT_URL")
	}
	
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       constants.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotifyAPI.AuthURL,
			TokenURL: spotifyAPI.TokenURL,
		},
	}, nil
}
