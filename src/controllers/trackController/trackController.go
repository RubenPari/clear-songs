package trackController

import (
	"sort"
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

	userTracks, errUserTracks := cacheManager.GetCachedUserTracksOrSet()

	if errUserTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting tracks",
		})
		return
	}

	// filter all tracks by artist
	tracksFilterers, errTracks := userService.GetAllUserTracksByArtist(idArtist, userTracks)

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

	// Explicitly invalidate cache after successful deletion
	cacheManager.InvalidateUserData()

	c.JSON(200, gin.H{
		"message": "Tracks deleted",
	})
}

// DeleteTrackByRange deletes tracks within a play count range from the user's library
func DeleteTrackByRange(c *gin.Context) {
	// get min query parameter (if exists)
	minStr := c.Query("min")
	minCount, _ := strconv.Atoi(minStr)

	maxStr := c.Query("max")
	maxCount, _ := strconv.Atoi(maxStr)

	userTracks, errUserTracks := cacheManager.GetCachedUserTracksOrSet()

	if errUserTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting tracks",
		})
		return
	}

	artistSummaryArray := trackHelper.GetArtistsSummary(userTracks)
	artistSummaryFiltered := utils.FilterSummaryByRange(artistSummaryArray, minCount, maxCount)

	// Track if any deletions occurred
	deletionsOccurred := false

	// delete all tracks from artists present in the summary object
	for _, artistObj := range artistSummaryFiltered {
		tracksFilters, errTracks := userService.GetAllUserTracksByArtist(spotifyAPI.ID(artistObj.Id), userTracks)

		if errTracks != nil {
			c.JSON(500, gin.H{
				"message": "Error getting tracks",
			})
			return
		}

		if len(tracksFilters) > 0 {
			// delete tracks from artist
			errDelete := userService.DeleteTracksUser(c, tracksFilters)

			if errDelete != nil {
				c.JSON(500, gin.H{
					"message": "Error deleting tracks",
				})
				return
			}
			deletionsOccurred = true
		}
	}

	// Only invalidate cache if deletions actually occurred
	if deletionsOccurred {
		cacheManager.InvalidateUserData()
	}

	c.JSON(200, gin.H{
		"message": "Tracks deleted",
	})
}

// GetTrackSummary returns a summary of tracks per artist, sorted by track count
func GetTrackSummary(c *gin.Context) {
	minStr := c.Query("min")
	maxStr := c.Query("max")
	minCount, _ := strconv.Atoi(minStr)
	maxCount, _ := strconv.Atoi(maxStr)

	tracks, errUserTracks := cacheManager.GetCachedUserTracksOrSet()

	if errUserTracks != nil {
		c.JSON(500, gin.H{
			"message": "Error getting tracks",
		})
		return
	}

	artistSummaryArray := trackHelper.GetArtistsSummary(tracks)

	// Apply range filters if provided
	if minStr != "" || maxStr != "" {
		artistSummaryArray = utils.FilterSummaryByRange(artistSummaryArray, minCount, maxCount)
	}

	// Sort by track count descending
	sort.Slice(artistSummaryArray, func(i, j int) bool {
		return artistSummaryArray[i].Count > artistSummaryArray[j].Count
	})

	c.JSON(200, artistSummaryArray)
}
