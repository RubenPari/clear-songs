package playlistController

import (
	"fmt"

	"github.com/RubenPari/clear-songs/src/services/playlistService"
	"github.com/RubenPari/clear-songs/src/services/userService"

	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetUserPlaylists returns all playlists owned or followed by the user.
//
// Returns a JSON array of playlists with id, name, and image_url.
func GetUserPlaylists(c *gin.Context) {
	playlists, err := playlistService.GetAllUserPlaylists()

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error getting user playlists",
			"error":   err.Error(),
		})
		return
	}

	// Convert to simplified format with only needed fields
	type PlaylistResponse struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		ImageURL string `json:"image_url,omitempty"`
	}

	var response []PlaylistResponse
	for _, playlist := range playlists {
		imageURL := ""
		if len(playlist.Images) > 0 {
			// Use the smallest image (usually the last one) for better performance
			for i := len(playlist.Images) - 1; i >= 0; i-- {
				if playlist.Images[i].Width <= 300 || i == 0 {
					imageURL = playlist.Images[i].URL
					break
				}
			}
		}

		response = append(response, PlaylistResponse{
			ID:       playlist.ID.String(),
			Name:     playlist.Name,
			ImageURL: imageURL,
		})
	}

	c.JSON(200, response)
}

// DeleteAllPlaylistTracks deletes all tracks from a playlist.
//
// The playlist ID is required and must be passed as a query parameter.
//
// The function first retrieves all tracks from the playlist, either from the cache
// or by calling the Spotify Web API. Then, it deletes all tracks from the playlist.
//
// If an error occurs while deleting the tracks, the function returns a JSON response
// with a 500 status code and an error message.
//
// Otherwise, it returns a JSON response with a 200 status code and a success message.
func DeleteAllPlaylistTracks(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Playlist id is required",
		})
		return
	}

	playlistID := spotifyAPI.ID(id)
	playlistTracks, errPlaylistTracks := cacheManager.GetCachedPlaylistTracksOrSet(playlistID)

	if errPlaylistTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting playlist tracks",
		})
		return
	}

	errDeletePlaylistTracks := playlistService.DeletePlaylistTracks(playlistID, playlistTracks)

	if errDeletePlaylistTracks != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error deleting tracks from playlist: id %s", id),
			"error":   errDeletePlaylistTracks.Error(),
		})
		return
	}

	// Explicitly invalidate cache for this specific playlist after successful deletion
	cacheManager.InvalidatePlaylist(playlistID)

	c.JSON(200, gin.H{
		"message": "Tracks deleted",
	})
}

// DeleteAllPlaylistAndUserTracks deletes all tracks
// from a playlist and from the user's library
func DeleteAllPlaylistAndUserTracks(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Playlist id is required",
		})
		return
	}

	playlistID := spotifyAPI.ID(id)
	playlistTracks, errPlaylistTracks := cacheManager.GetCachedPlaylistTracksOrSet(playlistID)

	if errPlaylistTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting playlist tracks",
		})
		return
	}

	// Save tracks backup before deletion
	errSaveTracksFile := utils.SaveTracksBackup(playlistTracks)

	if errSaveTracksFile != nil {
		c.JSON(500, gin.H{
			"message": "Error saving backup tracks to DB",
			"error":   errSaveTracksFile.Error(),
		})
		return
	}

	// Delete tracks from playlist
	errDeletePlaylistTracks := playlistService.DeletePlaylistTracks(playlistID, playlistTracks)

	if errDeletePlaylistTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error deleting tracks from playlist",
			"error":   errDeletePlaylistTracks.Error(),
		})
		return
	}

	// Convert playlist tracks to IDs for user library deletion
	playlistTracksIDs, errConvertIDs := utils.ConvertTracksToID(playlistTracks)

	if errConvertIDs != nil {
		c.JSON(500, gin.H{
			"message": "Error converting tracks to IDs",
			"error":   errConvertIDs.Error(),
		})
		return
	}

	// Delete tracks from user library
	errDeleteTrackUser := userService.DeleteTracksUser(c, playlistTracksIDs)

	if errDeleteTrackUser != nil {
		c.JSON(500, gin.H{
			"message": "Error deleting tracks from user library",
			"error":   errDeleteTrackUser.Error(),
		})
		return
	}

	// Explicitly invalidate cache after successful deletion
	// This operation affects both playlist and user data
	cacheManager.InvalidatePlaylist(playlistID)
	cacheManager.InvalidateUserData()

	c.JSON(200, gin.H{
		"message": "Tracks deleted",
	})
}
