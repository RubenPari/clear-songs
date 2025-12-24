package track

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/RubenPari/clear-songs/src/domain/entities"
	"github.com/RubenPari/clear-songs/src/domain/interfaces"
	"github.com/RubenPari/clear-songs/src/domain/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetTrackSummaryUseCase handles the business logic for getting track summaries
type GetTrackSummaryUseCase struct {
	spotifyRepo interfaces.SpotifyRepository
	cacheRepo  interfaces.CacheRepository
}

// NewGetTrackSummaryUseCase creates a new GetTrackSummaryUseCase
func NewGetTrackSummaryUseCase(
	spotifyRepo interfaces.SpotifyRepository,
	cacheRepo interfaces.CacheRepository,
) *GetTrackSummaryUseCase {
	return &GetTrackSummaryUseCase{
		spotifyRepo: spotifyRepo,
		cacheRepo:  cacheRepo,
	}
}

// Execute retrieves track summary grouped by artist, optionally filtered by range
func (uc *GetTrackSummaryUseCase) Execute(ctx context.Context, min, max int) ([]entities.ArtistSummary, error) {
	// 1. Check cache (if available)
	if uc.cacheRepo != nil {
		cacheKey := "track_summary"
		if min > 0 || max > 0 {
			cacheKey = fmt.Sprintf("track_summary_%d_%d", min, max)
		}
		
		var cached []entities.ArtistSummary
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
	summary := uc.calculateSummary(ctx, tracks, min, max)
	
	// 4. Sort by count descending
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].Count > summary[j].Count
	})
	
	// 5. Cache the result (if cache is available)
	if uc.cacheRepo != nil {
		cacheKey := "track_summary"
		if min > 0 || max > 0 {
			cacheKey = fmt.Sprintf("track_summary_%d_%d", min, max)
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

// calculateSummary calculates artist summary from tracks
func (uc *GetTrackSummaryUseCase) calculateSummary(
	ctx context.Context,
	tracks []spotifyAPI.SavedTrack,
	min, max int,
) []entities.ArtistSummary {
	// Group tracks by artist
	artistMap := make(map[string]struct {
		count int
		id    string
	})
	
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
			artistMap[artistName] = struct {
				count int
				id    string
			}{
				count: 1,
				id:    artistID,
			}
		}
	}
	
	// Convert to ArtistSummary array
	var summary []entities.ArtistSummary
	for artistName, data := range artistMap {
		// Apply range filter
		if min > 0 && data.count < min {
			continue
		}
		if max > 0 && data.count > max {
			continue
		}
		
		// Get artist image
		imageURL := ""
		if data.id != "" {
			artist, err := uc.spotifyRepo.GetArtist(ctx, spotifyAPI.ID(data.id))
			if err == nil && artist != nil {
				imageURL = utils.GetMediumImage(artist.Images)
			}
		}
		
		summary = append(summary, entities.ArtistSummary{
			ID:       data.id,
			Name:     artistName,
			Count:    data.count,
			ImageURL: imageURL,
		})
	}
	
	return summary
}
