package controllers

import (
	"github.com/RubenPari/clear-songs/src/services"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeleteAllPlaylistTracks deletes
// all tracks from a playlist
func DeleteAllPlaylistTracks(c *gin.Context) {
	idPlaylist := c.Query("id_playlist")

	if idPlaylist == "" {
		c.JSON(400, gin.H{
			"message": "Playlist id is required",
		})
		return
	}

	tracks, errTracks := services.GetAllPlaylistTracks(spotifyAPI.ID(idPlaylist))

	if errTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting tracks",
		})
		return
	}

	errDelete := services.DeleteTracksPlaylist(spotifyAPI.ID(idPlaylist), tracks)

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
	idPlaylist := c.Query("id_playlist")

	if idPlaylist == "" {
		c.JSON(400, gin.H{
			"message": "Playlist id is required",
		})
		return
	}

	tracks, errTracks := services.GetAllPlaylistTracks(spotifyAPI.ID(idPlaylist))

	if errTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting tracks",
		})
		return
	}

	errDeletePlaylist := services.DeleteTracksPlaylist(spotifyAPI.ID(idPlaylist), tracks)

	tracksIDs, errConvertIDs := utils.ConvertTracksToID(tracks)
	if errConvertIDs != nil {
		c.JSON(500, gin.H{
			"message": "Error converting tracks",
		})
		return
	}

	errDeleteUser := services.DeleteTracksUser(tracksIDs)

	if errDeletePlaylist != nil || errDeleteUser != nil {
		c.JSON(500, gin.H{
			"message": "Error deleting tracks",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Tracks deleted",
	})
}
