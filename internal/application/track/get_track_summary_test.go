package track

import (
	"context"
	"testing"

	"github.com/RubenPari/clear-songs/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	spotifyAPI "github.com/zmb3/spotify"
)

func TestGetTrackSummaryUseCase_Execute(t *testing.T) {
	// Setup mocks
	mockSpotifyRepo := new(mocks.MockSpotifyRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	
	useCase := NewGetTrackSummaryUseCase(mockSpotifyRepo, mockCacheRepo)
	ctx := context.Background()

	t.Run("Success - should return grouped summary", func(t *testing.T) {
		// Mock data
		tracks := []spotifyAPI.SavedTrack{
			{
				FullTrack: spotifyAPI.FullTrack{
					SimpleTrack: spotifyAPI.SimpleTrack{
						Name: "Song A",
						Artists: []spotifyAPI.SimpleArtist{
							{Name: "Artist 1", ID: "1"},
						},
					},
				},
			},
			{
				FullTrack: spotifyAPI.FullTrack{
					SimpleTrack: spotifyAPI.SimpleTrack{
						Name: "Song B",
						Artists: []spotifyAPI.SimpleArtist{
							{Name: "Artist 1", ID: "1"},
						},
					},
				},
			},
			{
				FullTrack: spotifyAPI.FullTrack{
					SimpleTrack: spotifyAPI.SimpleTrack{
						Name: "Song C",
						Artists: []spotifyAPI.SimpleArtist{
							{Name: "Artist 2", ID: "2"},
						},
					},
				},
			},
		}

		// Configure mock expectations
		mockCacheRepo.On("Get", ctx, "track_summary", mock.Anything).Return(false, nil)
		mockCacheRepo.On("GetUserTracks", ctx).Return(nil, nil)
		mockSpotifyRepo.On("GetAllUserTracks", ctx).Return(tracks, nil)
		mockCacheRepo.On("SetUserTracks", ctx, tracks, mock.Anything).Return(nil)
		
		// Mock Artist Image calls
		mockSpotifyRepo.On("GetArtist", ctx, spotifyAPI.ID("1")).Return(&spotifyAPI.FullArtist{}, nil)
		mockSpotifyRepo.On("GetArtist", ctx, spotifyAPI.ID("2")).Return(&spotifyAPI.FullArtist{}, nil)
		
		mockCacheRepo.On("Set", ctx, "track_summary", mock.Anything, mock.Anything).Return(nil)

		// Execute
		result, err := useCase.Execute(ctx, 0, 0)

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		
		// Check Artist 1 (should have 2 tracks)
		assert.Equal(t, "Artist 1", result[0].Name)
		assert.Equal(t, 2, result[0].Count)
		
		// Check Artist 2 (should have 1 track)
		assert.Equal(t, "Artist 2", result[1].Name)
		assert.Equal(t, 1, result[1].Count)
	})
}
