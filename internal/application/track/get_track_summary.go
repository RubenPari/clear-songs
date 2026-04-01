package track

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/RubenPari/clear-songs/internal/domain/track"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetTrackSummaryUseCase handles the business logic for getting track summaries
type GetTrackSummaryUseCase struct {
	spotifyRepo shared.SpotifyRepository
	cacheRepo   shared.CacheRepository
}

// NewGetTrackSummaryUseCase creates a new GetTrackSummaryUseCase
func NewGetTrackSummaryUseCase(
	spotifyRepo shared.SpotifyRepository,
	cacheRepo shared.CacheRepository,
) *GetTrackSummaryUseCase {
	return &GetTrackSummaryUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:   cacheRepo,
	}
}

// Execute retrieves track summary grouped by artist, optionally filtered by range and genre
func (uc *GetTrackSummaryUseCase) Execute(ctx context.Context, min, max int, genre string) ([]track.ArtistSummary, error) {
	// 1. Check cache (if available)
	if uc.cacheRepo != nil {
		cacheKey := "track_summary"
		if min > 0 || max > 0 {
			cacheKey = fmt.Sprintf("track_summary_%d_%d_%s", min, max, genre)
		}

		var cached []track.ArtistSummary
		if found, _ := uc.cacheRepo.Get(ctx, cacheKey, &cached); found {
			return cached, nil
		}
	}

	// 2. Get user tracks (from cache or API)
	tracks, err := uc.getUserTracks(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Calculate summary
	summary := uc.calculateSummary(ctx, tracks, min, max, genre)

	// 4. Sort by count descending
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].Count > summary[j].Count
	})

	// 5. Cache the result (if cache is available)
	if uc.cacheRepo != nil {
		cacheKey := "track_summary"
		if min > 0 || max > 0 {
			cacheKey = fmt.Sprintf("track_summary_%d_%d_%s", min, max, genre)
		}
		_ = uc.cacheRepo.Set(ctx, cacheKey, summary, 5*time.Minute)
	}

	return summary, nil
}

// getUserTracks retrieves tracks from cache or API
func (uc *GetTrackSummaryUseCase) getUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	// Try cache first (if available)
	if uc.cacheRepo != nil {
		cached, err := uc.cacheRepo.GetUserTracks(ctx)
		if err == nil && cached != nil && len(cached) > 0 {
			return cached, nil
		}
	}

	// Fetch from API
	tracks, err := uc.spotifyRepo.GetAllUserTracks(ctx)
	if err != nil {
		return nil, err
	}

	// Cache for future use (if cache is available)
	if uc.cacheRepo != nil {
		_ = uc.cacheRepo.SetUserTracks(ctx, tracks, 5*time.Minute)
	}

	return tracks, nil
}

// artistData holds artist information for summary calculation
type artistData struct {
	count  int
	id     string
	genres []string
}

// calculateSummary calculates artist summary from tracks
func (uc *GetTrackSummaryUseCase) calculateSummary(
	ctx context.Context,
	tracks []spotifyAPI.SavedTrack,
	min, max int,
	genre string,
) []track.ArtistSummary {
	// Group tracks by artist
	artistMap := make(map[string]artistData)

	for _, track := range tracks {
		if len(track.Artists) == 0 {
			continue
		}

		artistName := track.Artists[0].Name
		artistID := string(track.Artists[0].ID)

		if existing, exists := artistMap[artistName]; exists {
			existing.count++
			artistMap[artistName] = existing
		} else {
			artistMap[artistName] = artistData{
				count: 1,
				id:    artistID,
			}
		}
	}

	// Convert to ArtistSummary array
	var summary []track.ArtistSummary
	for artistName, data := range artistMap {
		// Apply range filter
		if min > 0 && data.count < min {
			continue
		}
		if max > 0 && data.count > max {
			continue
		}

		// Get artist image and genres
		imageURL := ""
		var genres []string

		if data.id != "" {
			artist, err := uc.spotifyRepo.GetArtist(ctx, spotifyAPI.ID(data.id))
			if err == nil && artist != nil {
				imageURL = utils.GetMediumImage(artist.Images)
				genres = artist.Genres
			}
		}

		// Apply genre filter if specified
		if genre != "" {
			genreLower := strings.ToLower(genre)
			found := false
			for _, g := range genres {
				if strings.Contains(strings.ToLower(g), genreLower) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Ensure genres is not nil for JSON marshaling
		if genres == nil {
			genres = []string{}
		}

		summary = append(summary, track.ArtistSummary{
			ID:       data.id,
			Name:     artistName,
			Count:    data.count,
			ImageURL: imageURL,
			Genres:   genres,
		})
	}

	return summary
}
