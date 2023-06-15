package controllers

import (
	"strconv"

	"github.com/RubenPari/clear-songs/src/lib/artist"
	"github.com/RubenPari/clear-songs/src/lib/user"
	"github.com/RubenPari/clear-songs/src/lib/utils"

	artistSpotifyLib "github.com/RubenPari/clear-songs/src/lib/artist"
	userSpotifyLib "github.com/RubenPari/clear-songs/src/lib/user"
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

	// get file query parameter (if exists)
	file := c.Query("file")

	// get tracks from user
	tracks, errTracks := userSpotifyLib.GetAllUserTracks()

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// get artist summary array
	artistSummaryArray := artist.GetArtistsSummary(tracks)

	// filter artist summary by min and max, if exists
	artistSummaryFiltered := utils.FilterSummaryByRange(artistSummaryArray, min, max)

	// if file query parameter exists, save a file with the summary

	// minimal -> only artists name
	if file == "minimal" {
		errFile := utils.SaveSummaryToFile(artistSummaryFiltered, true)

		if errFile != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Error saving file",
			})
			return
		}
	}

	// complete -> artists name and tracks
	if file == "complete" {
		errFile := utils.SaveSummaryToFile(artistSummaryFiltered, false)

		if errFile != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Error saving file",
			})
			return
		}
	}

	c.JSON(200, artistSummaryFiltered)
}

// DeleteTrackByArtist deletes all tracks from an artist
func DeleteTrackByArtist(c *gin.Context) {
	// get artist id from url
	idArtistString := c.Param("id_artist")
	idArtist := spotifyAPI.ID(idArtistString)

	// get tracks from user by artist
	tracksFilterers, errTracks := userSpotifyLib.GetAllUserTracksByArtist(idArtist)

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// delete tracks from artist
	errDelete := userSpotifyLib.DeleteTracksUser(tracksFilterers)

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
	tracksFilterers, errTracks := userSpotifyLib.GetAllUserTracksByGenre(name)

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// delete tracks from artist
	errDelete := userSpotifyLib.DeleteTracksUser(tracksFilterers)

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
	tracks, errTracks := userSpotifyLib.GetAllUserTracks()

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
	artistSummaryFiltered := utils.FilterSummaryByRange(artistSummaryArray, min, max)

	// delete all tracks from artists present
	// in the summary object
	for artistObj := range artistSummaryFiltered {
		// get tracks from user by artist
		tracksFilters, errTracks := user.GetAllUserTracksByArtist(spotifyAPI.ID(rune(artistObj)))

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

// DeleteTrackByFile deletes tracks from a file
// of tipe .txt with every artist name in a
// new line, send the file in the body of the
// request
func DeleteTrackByFile(c *gin.Context) {
	// get file from request
	file, errFile := c.FormFile("file")

	if errFile != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting file",
		})
		return
	}

	// get artists from file
	artists, errArtists := artistSpotifyLib.GetArtistsFromFile(file)

	if errArtists != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting artists",
		})
		return
	}

	// delete tracks from artists
	errDelete := userSpotifyLib.DeleteTracksByArtists(artists)

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
