package handlers

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/application/track"
	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// TrackControllerComplete is the complete refactored track controller
type TrackControllerComplete struct {
	BaseController
	getTrackSummaryUseCase *track.GetTrackSummaryUseCase
	deleteTracksByArtistUC *track.DeleteTracksByArtistUseCase
	deleteTracksByRangeUC  *track.DeleteTracksByRangeUseCase
	deleteTrackUC          *track.DeleteTrackUseCase
	getTracksByArtistUC    *track.GetTracksByArtistUseCase
}

// NewTrackControllerComplete creates a new complete track controller
func NewTrackControllerComplete(
	getTrackSummaryUC *track.GetTrackSummaryUseCase,
	deleteByArtistUC *track.DeleteTracksByArtistUseCase,
	deleteByRangeUC *track.DeleteTracksByRangeUseCase,
	getTracksByArtistUC *track.GetTracksByArtistUseCase,
	deleteTrackUC *track.DeleteTrackUseCase,
) *TrackControllerComplete {
	return &TrackControllerComplete{
		getTrackSummaryUseCase: getTrackSummaryUC,
		deleteTracksByArtistUC: deleteByArtistUC,
		deleteTracksByRangeUC:  deleteByRangeUC,
		deleteTrackUC:          deleteTrackUC,
		getTracksByArtistUC:    getTracksByArtistUC,
	}
}

// GetTrackSummary handles GET /track/summary
func (tc *TrackControllerComplete) GetTrackSummary(c *gin.Context) {
	var req track.RangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tc.JSONValidationError(c, "Invalid min or max parameters")
		return
	}

	// Execute use case
	// Note: the original manual validation fell back to 0 if min/max strings were empty,
	// which matches how Gin parses missing query integers.
	ctx := context.Background() // In production, c.Request.Context() is preferred.
	result, err := tc.getTrackSummaryUseCase.Execute(ctx, req.Min, req.Max)
	if err != nil {
		tc.HandleDomainError(c, err)
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
	ctx := context.Background() // In production, use c.Request.Context()
	tracks, err := tc.getTracksByArtistUC.Execute(ctx, artistID)
	if err != nil {
		tc.HandleDomainError(c, err)
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
	ctx := context.Background() // In production, use c.Request.Context()
	if err := tc.deleteTracksByArtistUC.Execute(ctx, artistID); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}

// DeleteTrack handles DELETE /track/:id_track
func (tc *TrackControllerComplete) DeleteTrack(c *gin.Context) {
	// Get track ID from URL
	idTrackString := c.Param("id_track")
	if idTrackString == "" {
		tc.JSONValidationError(c, "Track ID is required")
		return
	}

	trackID := spotifyAPI.ID(idTrackString)

	// Execute use case
	ctx := context.Background() // In production, use c.Request.Context()
	if err := tc.deleteTrackUC.Execute(ctx, trackID); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Track deleted successfully"})
}

// DeleteTrackByRange handles DELETE /track/by-range
func (tc *TrackControllerComplete) DeleteTrackByRange(c *gin.Context) {
	var req track.RangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tc.JSONValidationError(c, "Invalid min or max parameters")
		return
	}

	// At least one parameter must be provided for a destructive action
	if req.Min == 0 && req.Max == 0 && c.Query("min") == "" && c.Query("max") == "" {
		tc.JSONValidationError(c, "At least one of min or max must be provided")
		return
	}

	// Execute use case
	ctx := context.Background() // In production, use c.Request.Context()
	if err := tc.deleteTracksByRangeUC.Execute(ctx, req.Min, req.Max); err != nil {
		tc.HandleDomainError(c, err)
		return
	}

	tc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}
