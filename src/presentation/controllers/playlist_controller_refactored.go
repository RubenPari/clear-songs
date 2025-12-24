package controllers

import (
	"context"

	"github.com/RubenPari/clear-songs/src/application/usecases/playlist"
	"github.com/RubenPari/clear-songs/src/domain/utils"
	"github.com/RubenPari/clear-songs/src/presentation/dto"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// PlaylistControllerRefactored is the refactored playlist controller using dependency injection
type PlaylistControllerRefactored struct {
	BaseController
	getUserPlaylistsUC          *playlist.GetUserPlaylistsUseCase
	deletePlaylistTracksUC      *playlist.DeletePlaylistTracksUseCase
	deletePlaylistAndLibraryUC  *playlist.DeletePlaylistAndLibraryTracksUseCase
}

// NewPlaylistControllerRefactored creates a new playlist controller
func NewPlaylistControllerRefactored(
	getUserPlaylistsUC *playlist.GetUserPlaylistsUseCase,
	deletePlaylistTracksUC *playlist.DeletePlaylistTracksUseCase,
	deletePlaylistAndLibraryUC *playlist.DeletePlaylistAndLibraryTracksUseCase,
) *PlaylistControllerRefactored {
	return &PlaylistControllerRefactored{
		getUserPlaylistsUC:         getUserPlaylistsUC,
		deletePlaylistTracksUC:     deletePlaylistTracksUC,
		deletePlaylistAndLibraryUC: deletePlaylistAndLibraryUC,
	}
}

// GetUserPlaylists handles GET /playlist/list
func (pc *PlaylistControllerRefactored) GetUserPlaylists(c *gin.Context) {
	ctx := context.Background()
	playlists, err := pc.getUserPlaylistsUC.Execute(ctx)
	if err != nil {
		pc.JSONInternalError(c, "Error getting user playlists")
		return
	}

	// Convert to response format
	var response []dto.PlaylistResponse
	for _, playlist := range playlists {
		imageURL := utils.GetMediumImage(playlist.Images)

		response = append(response, dto.PlaylistResponse{
			ID:       playlist.ID.String(),
			Name:     playlist.Name,
			ImageURL: imageURL,
		})
	}

	pc.JSONSuccess(c, response)
}

// DeleteAllPlaylistTracks handles DELETE /playlist/delete-tracks
func (pc *PlaylistControllerRefactored) DeleteAllPlaylistTracks(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		pc.JSONValidationError(c, "Playlist id is required")
		return
	}

	playlistID := spotifyAPI.ID(id)
	ctx := context.Background()

	if err := pc.deletePlaylistTracksUC.Execute(ctx, playlistID); err != nil {
		pc.JSONInternalError(c, "Error deleting tracks from playlist")
		return
	}

	pc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}

// DeleteAllPlaylistAndUserTracks handles DELETE /playlist/delete-tracks-and-library
func (pc *PlaylistControllerRefactored) DeleteAllPlaylistAndUserTracks(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		pc.JSONValidationError(c, "Playlist id is required")
		return
	}

	playlistID := spotifyAPI.ID(id)
	ctx := context.Background()

	if err := pc.deletePlaylistAndLibraryUC.Execute(ctx, playlistID); err != nil {
		pc.JSONInternalError(c, "Error deleting tracks from playlist and library")
		return
	}

	pc.JSONSuccess(c, gin.H{"message": "Tracks deleted successfully"})
}
