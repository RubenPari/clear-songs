package di

import (
	"errors"
	"log"
	"os"

	"github.com/RubenPari/clear-songs/internal/application/auth"
	"github.com/RubenPari/clear-songs/internal/application/playlist"
	"github.com/RubenPari/clear-songs/internal/application/track"
	"github.com/RubenPari/clear-songs/internal/domain/shared/constants"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/redis"
	"github.com/RubenPari/clear-songs/internal/infrastructure/external/spotify"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// Container holds all application dependencies
type Container struct {
	// Repositories (as interfaces)
	SpotifyRepo  shared.SpotifyRepository
	CacheRepo    shared.CacheRepository
	DatabaseRepo shared.DatabaseRepository
	
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
	var cacheRepo shared.CacheRepository
	redisCache, err := redis.NewRedisCacheRepository()
	if err != nil {
		log.Printf("WARNING: Cache repository initialization failed: %v", err)
		log.Println("WARNING: Application will continue without Redis caching")
		// Use no-op cache implementation
		cacheRepo = redis.NewNoOpCacheRepository()
	} else {
		cacheRepo = redisCache
	}
	
	// Initialize auth use cases
	loginUC := auth.NewLoginUseCase(oauthConfig)
	callbackUC := auth.NewCallbackUseCase(oauthConfig, spotifyRepo, cacheRepo)
	logoutUC := auth.NewLogoutUseCase(spotifyRepo, cacheRepo)
	isAuthUC := auth.NewIsAuthUseCase(spotifyRepo)
	
	// Initialize database repository (may be nil if database not available)
	databaseRepo := postgres.NewPostgresRepository(postgres.Db)
	
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
