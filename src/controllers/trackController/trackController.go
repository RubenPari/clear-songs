package trackController

import (
	"strconv"

	"github.com/RubenPari/clear-songs/src/helpers/trackHelper"
	"github.com/RubenPari/clear-songs/src/services/userService"

	cacheManager "github.com/RubenPari/clear-songs/src/cache"

	"github.com/RubenPari/clear-songs/src/utils"

	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// DeleteTrackByArtist deletes all tracks from an artist
func DeleteTrackByArtist(c *gin.Context) {
	// get artist id from url
	idArtistString := c.Param("id_artist")
	idArtist := spotifyAPI.ID(idArtistString)

	// get all tracks from user
	var tracks []spotifyAPI.SavedTrack
	var errTracks error

	value, found := cacheManager.Get("userTracks")

	if found {
		tracks = value.([]spotifyAPI.SavedTrack)
	} else {
		tracks, errTracks = userService.GetAllUserTracks(c)

		if errTracks != nil {
			c.JSON(500, gin.H{
				"message": "Error getting tracks",
			})
			return
		}

		// save user tracks in cacheManager
		cacheManager.Set("userTracks", tracks)
	}

	// filter all tracks by artist
	tracksFilterers, errTracks := userService.GetAllUserTracksByArtist(idArtist, tracks)

	if errTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting tracks",
		})
		return
	}

	// delete tracks from artist
	errDelete := userService.DeleteTracksUser(c, tracksFilterers)

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

func DeleteTrackByRange(c *gin.Context) {
	// get min query parameter (if exists)
	minStr := c.Query("min")
	minCount, _ := strconv.Atoi(minStr)

	maxStr := c.Query("max")
	maxCount, _ := strconv.Atoi(maxStr)

	// get tracks from user
	var tracks []spotifyAPI.SavedTrack
	var errTracks error

	value, found := cacheManager.Get("userTracks")

	if found {
		tracks = value.([]spotifyAPI.SavedTrack)
	} else {
		tracks, errTracks = userService.GetAllUserTracks(c)

		if errTracks != nil {
			c.JSON(500, gin.H{
				"message": "Error getting tracks",
			})
			return
		}

		// save user tracks in cacheManager
		cacheManager.Set("userTracks", tracks)
	}

	artistSummaryArray := trackHelper.GetArtistsSummary(tracks)

	artistSummaryFiltered := utils.FilterSummaryByRange(artistSummaryArray, minCount, maxCount)

	// delete all tracks from artists present
	// in the summary object
	for artistObj := range artistSummaryFiltered {
		tracksFilters, errTracks := userService.GetAllUserTracksByArtist(spotifyAPI.ID(rune(artistObj)), tracks)

		if errTracks != nil {
			c.JSON(500, gin.H{
				"message": "Error getting tracks",
			})
			return
		}

		// delete tracks from artist
		errDelete := userService.DeleteTracksUser(c, tracksFilters)

		if errDelete != nil {
			c.JSON(500, gin.H{
				"message": "Error deleting tracks",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "Tracks deleted",
	})
}
