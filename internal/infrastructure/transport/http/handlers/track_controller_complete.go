package handlers

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/application/track"
	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/RubenPari/clear-songs/internal/application/shared/validators"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// TrackControllerComplete is the complete refactored track controller
type TrackControllerComplete struct {
	BaseController
	getTrackSummaryUseCase    *track.GetTrackSummaryUseCase
	deleteTracksByArtistUC    *track.DeleteTracksByArtistUseCase
	deleteTracksByRangeUC     *track.DeleteTracksByRangeUseCase
	getTracksByArtistUC       *track.GetTracksByArtistUseCase
}

// NewTrackControllerComplete creates a new complete track controller
func NewTrackControllerComplete(
	getTrackSummaryUC *track.GetTrackSummaryUseCase,
	deleteByArtistUC *track.DeleteTracksByArtistUseCase,
	deleteByRangeUC *track.DeleteTracksByRangeUseCase,
	getTracksByArtistUC *track.GetTracksByArtistUseCase,
) *TrackControllerComplete {
	return &TrackControllerComplete{
		getTrackSummaryUseCase: getTrackSummaryUC,
		deleteTracksByArtistUC: deleteByArtistUC,
		deleteTracksByRangeUC: deleteByRangeUC,
		getTracksByArtistUC:   getTracksByArtistUC,
	}
}

// GetTrackSummary handles GET /track/summary
func (tc *TrackControllerComplete) GetTrackSummary(c *gin.Context) {
	// Validate and parse query parameters
	params, err := validators.ParseRangeParams(
		c.Query("min"),
		c.Query("max"),
	)
	if err != nil {
		tc.JSONValidationError(c, err.Error())
		return
	}

	min, max := params.GetMinMax()

	// Execute use case
	ctx := context.Background()
	result, err := tc.getTrackSummaryUseCase.Execute(ctx, min, max)
	if err != nil {
		tc.JSONInternalError(c, "Failed to get track summary")
		return
	}

	// Convert to API response format (entities to models)
	var response []track.ArtistSummary
	for _, artist := range result {
		response = append(response, track.ArtistSummary{
			Id:       artist.ID,
			Name:     artist.Name,
			Count:    artist.Count,
			ImageURL: artist.ImageURL,
		})
	}

	tc.JSONSuccess(c, response)
}

// GetTracksByArtist handles GET /track/by-artist/:id_artist
func (tc *TrackControllerComplete) GetTracksByArtist(c *gin.Context) {
	// Get artist ID from URL
	idArtistString := c.Param("id_artist")
	if idArtistString == "" {
		tc.JSONValidationError(c, "Artist ID is required")
		return
	}

	artistID := spotifyAPI.ID(idArtistString)

	// Execute use case
	ctx := context.Background()
	tracks, err := tc.getTracksByArtistUC.Execute(ctx, artistID)
	if err != nil {
		tc.JSONInternalError(c, "Failed to get tracks by artist")
		return
	}

	// Convert to response format
	var response []track.TrackResponse
	for _, t := range tracks {
		artists := make([]string, len(t.Artists))
		for i, artist := range t.Artists {
			artists[i] = artist.Name
		}

		imageURL := utils.GetMediumImage(t.Album.Images)

		spotifyURL := ""
		if url, exists := t.ExternalURLs["spotify"]; exists {
			spotifyURL = url
		}

		response = append(response, track.TrackResponse{
			ID:         t.ID.String(),
			Name:       t.Name,
			Artists:    artists,
			Album:      t.Album.Name,
			Duration:   t.Duration,
			ImageURL:   imageURL,
			SpotifyURL: spotifyURL,
		})
	}

	tc.JSONSuccess(c, response)
}

// DeleteTrackByArtist handles DELETE /track/by-artist/:id_artist
func (tc *TrackControllerComplete) DeleteTrackByArtist(c *gin.Context) {
	// Get artist ID from URL
	idArtistString := c.Param("id_artist")
	if idArtistString == "" {
		tc.JSONValidationError(c, "Artist ID is required")
		return
	}

	artistID := spotifyAPI.ID(idArtistString)

	// Execute use case
	ctx := context.Background()
	if err := tc.deleteTracksByArtistUC.Execute(ctx, artistID); err != nil {
		tc.JSONInternalError(c, "Failed to delete tracks by artist")
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}

// DeleteTrackByRange handles DELETE /track/by-range
func (tc *TrackControllerComplete) DeleteTrackByRange(c *gin.Context) {
	// Validate and parse query parameters
	params, err := validators.ParseRangeParams(
		c.Query("min"),
		c.Query("max"),
	)
	if err != nil {
		tc.JSONValidationError(c, err.Error())
		return
	}

	min, max := params.GetMinMax()

	// At least one parameter must be provided
	if params.Min == nil && params.Max == nil {
		tc.JSONValidationError(c, "At least one of min or max must be provided")
		return
	}

	// Execute use case
	ctx := context.Background()
	if err := tc.deleteTracksByRangeUC.Execute(ctx, min, max); err != nil {
		tc.JSONInternalError(c, "Failed to delete tracks by range")
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}
