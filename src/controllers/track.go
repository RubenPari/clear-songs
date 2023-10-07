package controllers

import (
	"strconv"

	"github.com/RubenPari/clear-songs/src/services/artist"
	artistService "github.com/RubenPari/clear-songs/src/services/artist"
	userService "github.com/RubenPari/clear-songs/src/services/user"
	"github.com/RubenPari/clear-songs/src/utils"

	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetTrackSummary returns a summary of the user's tracks
func GetTrackSummary(c *gin.Context) {
	// get min query parameter (if exists)
	minStr := c.Query("min")
	minCount, _ := strconv.Atoi(minStr)

	// get max query parameter (if exists)
	maxStr := c.Query("max")
	maxCount, _ := strconv.Atoi(maxStr)

	// get group-by query parameter (if exists)
	groupByGenreStr := c.Query("group-by")
	groupByGenre := groupByGenreStr == "1"

	// get tracks from user
	tracks, errTracks := userService.GetAllUserTracks()

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// get artist summary array
	artistSummaryArray := artist.GetArtistsSummary(tracks)

	// manage case when groupBy is present or not
	if groupByGenre {
		// group artist summary by genres
		artistSummaryArrayGrouped := artist.GroupArtistSummaryByGenres(artistSummaryArray)

		// filter artist summary group by min and max, if exists
		artistSummaryFiltered := utils.FilterGroupSummaryByRange(artistSummaryArrayGrouped, minCount, maxCount)

		c.JSON(200, artistSummaryFiltered)
	} else {
		// filter artist summary by min and max
		artistSummaryFiltered := utils.FilterSummaryByRange(artistSummaryArray, minCount, maxCount)

		c.JSON(200, artistSummaryFiltered)
	}
}

// DeleteTrackByArtist deletes all tracks from an artist
func DeleteTrackByArtist(c *gin.Context) {
	// get artist id from url
	idArtistString := c.Param("id_artist")
	idArtist := spotifyAPI.ID(idArtistString)

	// get tracks from user by artist
	tracksFilterers, errTracks := userService.GetAllUserTracksByArtist(idArtist)

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// delete tracks from artist
	errDelete := userService.DeleteTracksUser(tracksFilterers)

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
	tracksFilterers, errTracks := userService.GetAllUserTracksByGenre(name)

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// delete tracks from artist
	errDelete := userService.DeleteTracksUser(tracksFilterers)

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
	minCount, _ := strconv.Atoi(minStr)

	// get max query parameter (if exists)
	maxStr := c.Query("max")
	maxCount, _ := strconv.Atoi(maxStr)

	// get tracks from user
	tracks, errTracks := userService.GetAllUserTracks()

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
		return
	}

	// get artist summary array
	artistSummaryArray := artistService.GetArtistsSummary(tracks)

	// filter artist summary by min and max
	artistSummaryFiltered := utils.FilterSummaryByRange(artistSummaryArray, minCount, maxCount)

	// delete all tracks from artists present
	// in the summary object
	for artistObj := range artistSummaryFiltered {
		// get tracks from user by artist
		tracksFilters, errTracks := userService.GetAllUserTracksByArtist(spotifyAPI.ID(rune(artistObj)))

		if errTracks != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Error getting tracks",
			})
			return
		}

		// delete tracks from artist
		errDelete := userService.DeleteTracksUser(tracksFilters)

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
	artists, errArtists := artistService.GetArtistsFromFile(file)

	if errArtists != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting artists",
		})
		return
	}

	// delete tracks from artists
	errDelete := userService.DeleteTracksByArtists(artists)

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
