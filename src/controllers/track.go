package controllers

import (
	"github.com/RubenPari/clear-songs/src/lib/array"
	"github.com/RubenPari/clear-songs/src/lib/artist"
	"github.com/RubenPari/clear-songs/src/lib/user"
	"strconv"

	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetTrackSummary returns a summary of the user's tracks
func GetTrackSummary(c *gin.Context) {
	// get min query parameter (if exists)
	minStr := c.Query("min")
	min, _ := strconv.Atoi(minStr)

	// get max query parameter (if exists)
	maxStr := c.Query("max")
	max, _ := strconv.Atoi(maxStr)

	// get tracks from user
	tracks, errTracks := user.GetAllUserTracks()

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// get artist summary array
	artistSummaryArray := artist.GetArtistsSummary(tracks)

	// filter artist summary by min and max
	artistSummaryFiltered := array.FilterSummaryByRange(artistSummaryArray, min, max)

	c.JSON(200, artistSummaryFiltered)
}

// DeleteTrackByArtist deletes all tracks from an artist
func DeleteTrackByArtist(c *gin.Context) {
	// get artist id from url
	idArtistString := c.Param("id_artist")
	idArtist := spotifyAPI.ID(idArtistString)

	// get tracks from user by artist
	tracksFilterers, errTracks := user.GetAllUserTracksByArtist(idArtist)

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// delete tracks from artist
	errDelete := user.DeleteTracksUser(tracksFilterers)

	if errDelete != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting tracks",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Tracks deleted",
	})
}

func DeleteTrackByGenre(c *gin.Context) {
	// get genre name from query param
	name := c.Query("name")

	// get tracks from user
	tracksFilterers, errTracks := user.GetAllUserTracksByGenre(name)

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// delete tracks from artist
	errDelete := user.DeleteTracksUser(tracksFilterers)

	if errDelete != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting tracks",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Tracks deleted",
	})
}

func DeleteTrackByRange(c *gin.Context) {
	// get min query parameter (if exists)
	minStr := c.Query("min")
	min, _ := strconv.Atoi(minStr)

	// get max query parameter (if exists)
	maxStr := c.Query("max")
	max, _ := strconv.Atoi(maxStr)

	// get tracks from user
	tracks, errTracks := user.GetAllUserTracks()

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// get artist summary array
	artistSummaryArray := artist.GetArtistsSummary(tracks)

	// filter artist summary by min and max
	artistSummaryFiltered := array.FilterSummaryByRange(artistSummaryArray, min, max)

	// delete all tracks from artists present
	// in the summary object
	for artistObj := range artistSummaryFiltered {
		// get tracks from user by artist
		tracksFilters, errTracks := user.GetAllUserTracksByArtist(spotifyAPI.ID(artistObj))

		if errTracks != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Error getting tracks",
			})
			return
		}

		// delete tracks from artist
		errDelete := user.DeleteTracksUser(tracksFilters)

		if errDelete != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Error deleting tracks",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Tracks deleted",
	})
}
