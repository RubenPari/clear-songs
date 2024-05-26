package controllers

import (
	"github.com/RubenPari/clear-songs/src/services"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeleteAllPlaylistTracks deletes all tracks
// from a playlist and from the user's library
func DeleteAllPlaylistTracks(c *gin.Context) {
	idPlaylist := c.Param("id")

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
