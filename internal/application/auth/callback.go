package auth

import (
	"context"
	"os"

	domainAuth "github.com/RubenPari/clear-songs/internal/domain/auth"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"golang.org/x/oauth2"
)

// CallbackUseCase handles the business logic for OAuth callback
type CallbackUseCase struct {
	oauthConfig *oauth2.Config
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
	userRepo    domainAuth.UserRepository
}

// NewCallbackUseCase creates a new CallbackUseCase
func NewCallbackUseCase(
	oauthConfig *oauth2.Config,
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
	userRepo domainAuth.UserRepository,
) *CallbackUseCase {
	return &CallbackUseCase{
		oauthConfig: oauthConfig,
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
		userRepo:    userRepo,
	}
}

// Execute processes the OAuth callback and returns the frontend redirect URL
func (uc *CallbackUseCase) Execute(ctx context.Context, code string, localUserID string) (string, error) {
	// 1. Exchange code for token
	token, err := uc.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	// 2. Save token to cache
	if uc.cacheRepo != nil {
		if err := uc.cacheRepo.SetToken(ctx, token); err != nil {
			// Log error but continue
		}
	}

	// 3. Set token in Spotify repository
	if err := uc.spotifyRepo.SetAccessToken(token); err != nil {
		return "", err
	}

	// 4. Verify authentication by getting current user
	spotifyUser, err := uc.spotifyRepo.GetCurrentUser(ctx)
	if err != nil {
		return "", err
	}

	// Link Spotify profile to Local User if logged in
	if localUserID != "" && uc.userRepo != nil {
		localUser, err := uc.userRepo.GetByID(ctx, localUserID)
		if err == nil && localUser != nil {
			localUser.SpotifyID = &spotifyUser.ID
			_ = uc.userRepo.Update(ctx, localUser)
		}
	}

	// 5. Get frontend URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4200"
	}

	return frontendURL + "/callback", nil
}
