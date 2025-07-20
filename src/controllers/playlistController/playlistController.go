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

// DeleteAllPlaylistTracks godoc
// @Summary Delete all tracks from playlist
// @Schemes
// @Description Removes all tracks from a specified playlist
// @Tags playlist
// @Accept json
// @Produce json
// @Param id query string true "Playlist ID"
// @Success 200 {object} map[string]string "message: Tracks deleted"
// @Failure 400 {object} map[string]string "message: Playlist id is required"
// @Failure 500 {object} map[string]string "message: Error deleting tracks from playlist"
// @Router /playlist/tracks [delete]
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

// DeleteAllPlaylistAndUserTracks godoc
// @Summary Delete tracks from playlist and user library
// @Schemes
// @Description Removes all tracks from the playlist and user's library and save the tracks on DB for backup before deleting them
// @Tags playlist
// @Accept json
// @Produce json
// @Param id query string true "Playlist ID"
// @Success 200 {object} map[string]string "message: Tracks deleted"
// @Failure 400 {object} map[string]string "message: Playlist id is required"
// @Failure 500 {object} map[string]string "message: Error deleting tracks"
// @Router /playlist/tracks/all [delete]
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
