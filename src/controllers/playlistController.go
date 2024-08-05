package controllers

import (
	"github.com/RubenPari/clear-songs/src/cacheManager"
	"github.com/RubenPari/clear-songs/src/services"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeleteAllPlaylistTracks deletes
// all tracks from a playlist
func DeleteAllPlaylistTracks(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Playlist id is required",
		})
		return
	}

	// get all playlist tracks
	var tracksPlaylist []spotifyAPI.PlaylistTrack
	var errTrackPlaylist error

	value, found := cacheManager.Get("tracksPlaylist" + id)

	if found {
		tracksPlaylist = value.([]spotifyAPI.PlaylistTrack)
	} else {
		tracksPlaylist, errTrackPlaylist = services.GetAllPlaylistTracks(spotifyAPI.ID(id))

		if errTrackPlaylist != nil {
			c.JSON(500, gin.H{
				"message": "Error getting playlist tracks",
			})
			return
		}

		cacheManager.Set("tracksPlaylist"+id, tracksPlaylist)
	}

	errDelete := services.DeleteTracksPlaylist(spotifyAPI.ID(id), tracksPlaylist)

	if errDelete != nil {
		c.JSON(500, gin.H{
			"message": "Error deleting tracks",
		})
		return
	}

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

	// get all playlist tracks
	var tracksPlaylist []spotifyAPI.PlaylistTrack
	var errTrackPlaylist error

	value, _ := cacheManager.Get("tracksPlaylist" + id)
	if value != nil {
		tracksPlaylist = value.([]spotifyAPI.PlaylistTrack)
	} else {
		tracksPlaylist, errTrackPlaylist = services.GetAllPlaylistTracks(spotifyAPI.ID(id))

		if errTrackPlaylist != nil {
			c.JSON(500, gin.H{
				"message": "Error getting playlist tracks",
			})
			return
		}

		cacheManager.Set("tracksPlaylist"+id, tracksPlaylist)
	}

	errSaveTracksFile := utils.SaveTracksBackup(tracksPlaylist)

	if errSaveTracksFile != nil {
		c.JSON(500, gin.H{
			"message": "Error saving backup tracks to file",
		})
		return
	}

	errDeletePlaylistTracks := services.DeleteTracksPlaylist(spotifyAPI.ID(id), tracksPlaylist)

	tracksPlaylistIDs, errConvertIDs := utils.ConvertTracksToID(tracksPlaylist)

	if errConvertIDs != nil {
		c.JSON(500, gin.H{
			"message": "Error converting tracks to IDs",
		})
		return
	}

	errDeleteTrackUser := services.DeleteTracksUser(tracksPlaylistIDs)

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
