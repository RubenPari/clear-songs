package controllers

import (
	"fmt"
	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/services"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

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

	playlistTracks, errPlaylistTracks := cacheManager.GetCachedPlaylistTracksOrSet(spotifyAPI.ID(id))

	if errPlaylistTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting playlist tracks",
		})
		return
	}

	errDeletePlaylistTracks := services.DeletePlaylistTracks(spotifyAPI.ID(id), playlistTracks)

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
	var playlistTracks []spotifyAPI.PlaylistTrack
	var errPlaylistTracks error

	value, _ := cacheManager.Get("tracksPlaylist" + id)
	if value != nil {
		playlistTracks = value.([]spotifyAPI.PlaylistTrack)
	} else {
		playlistTracks, errPlaylistTracks = services.GetAllPlaylistTracks(spotifyAPI.ID(id))

		if errPlaylistTracks != nil {
			c.JSON(500, gin.H{
				"message": "Error getting playlist tracks",
			})
			return
		}

		cacheManager.Set("tracksPlaylist"+id, playlistTracks)
	}

	errSaveTracksFile := utils.SaveTracksBackup(playlistTracks)

	if errSaveTracksFile != nil {
		c.JSON(500, gin.H{
			"message": "Error saving backup tracks to DB",
			"error":   errSaveTracksFile.Error(),
		})
		return
	}

	errDeletePlaylistTracks := services.DeletePlaylistTracks(spotifyAPI.ID(id), playlistTracks)

	playlistTracksIDs, errConvertIDs := utils.ConvertTracksToID(playlistTracks)

	if errConvertIDs != nil {
		c.JSON(500, gin.H{
			"message": "Error converting tracks to IDs",
		})
		return
	}

	errDeleteTrackUser := services.DeleteTracksUser(playlistTracksIDs)

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

// CreatePlaylistTracksMinor create a playlist with all
// the tracks in the user library that belong to artists
// for which the user owns a maximum of 5 tracks
func CreatePlaylistTracksMinor(c *gin.Context) {
	// get all tracks from user
	var tracks []spotifyAPI.SavedTrack
	var errTracks error

	value, found := cacheManager.Get("userTracks")

	if found {
		tracks = value.([]spotifyAPI.SavedTrack)
	} else {
		tracks, errTracks = services.GetAllUserTracks()

		if errTracks != nil {
			c.JSON(500, gin.H{
				"message": "Error getting tracks",
			})
			return
		}

		// save user tracks in cacheManager
		cacheManager.Set("userTracks", tracks)
	}

	// create playlist
	errCreatePlaylistTracksMinor := services.CreatePlaylistTracksMinor(tracks)

	if errCreatePlaylistTracksMinor != nil {
		c.JSON(500, gin.H{
			"message": "Error creating playlist",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Playlist created",
	})
}
