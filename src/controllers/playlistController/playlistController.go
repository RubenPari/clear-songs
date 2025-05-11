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
	value := cacheManager.Get("modifiedCachedValue")

	if value == true {
		cacheManager.Reset()
	}

	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Playlist id is required",
		})
		return
	}

	playlistTracks, errPlaylistTracks := cacheManager.GetCachedPlaylistTracksOrSet(spotifyAPI.ID(id))

	if errPlaylistTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting playlist tracks",
		})
		return
	}

	errDeletePlaylistTracks := playlistService.DeletePlaylistTracks(spotifyAPI.ID(id), playlistTracks)

	if errDeletePlaylistTracks != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error deleting tracks from playlist: id %s", id),
			"error":   errDeletePlaylistTracks.Error(),
		})
		return
	}

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
	value := cacheManager.Get("modifiedCachedValue")

	if value == true {
		cacheManager.Reset()
	}

	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Playlist id is required",
		})
		return
	}

	playlistTracks, errPlaylistTracks := cacheManager.GetCachedPlaylistTracksOrSet(spotifyAPI.ID(id))

	if errPlaylistTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting playlist tracks",
		})
		return
	}

	errSaveTracksFile := utils.SaveTracksBackup(playlistTracks)

	if errSaveTracksFile != nil {
		c.JSON(500, gin.H{
			"message": "Error saving backup tracks to DB",
			"error":   errSaveTracksFile.Error(),
		})
		return
	}

	errDeletePlaylistTracks := playlistService.DeletePlaylistTracks(spotifyAPI.ID(id), playlistTracks)

	playlistTracksIDs, errConvertIDs := utils.ConvertTracksToID(playlistTracks)

	if errConvertIDs != nil {
		c.JSON(500, gin.H{
			"message": "Error converting tracks to IDs",
		})
		return
	}

	errDeleteTrackUser := userService.DeleteTracksUser(c, playlistTracksIDs)

	if errDeletePlaylistTracks != nil || errDeleteTrackUser != nil {
		c.JSON(500, gin.H{
			"message": "Error deleting tracks",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Tracks deleted",
	})
}
