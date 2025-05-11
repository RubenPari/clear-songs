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

// DeleteTrackByArtist godoc
// @Summary Delete all tracks by artist
// @Schemes
// @Description Removes all tracks from a specific artist from user's library
// @Tags track
// @Accept json
// @Produce json
// @Param id_artist path string true "Artist ID"
// @Success 200 {object} map[string]string "message: Tracks deleted"
// @Failure 500 {object} map[string]string "message: Error deleting tracks"
// @Router /track/artist/{id_artist} [delete]
// DeleteTrackByArtist deletes all tracks from an artist
func DeleteTrackByArtist(c *gin.Context) {
	spotifyClient := c.MustGet("spotifyClient").(*spotifyAPI.Client)

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

	c.JSON(200, gin.H{
		"message": "Tracks deleted",
	})
}

// DeleteTrackByRange godoc
// @Summary Delete tracks within play count range
// @Schemes
// @Description Removes tracks that fall within a specified play count range
// @Tags track
// @Accept json
// @Produce json
// @Param min query integer false "Minimum play count"
// @Param max query integer false "Maximum play count"
// @Success 200 {object} map[string]string "message: Tracks deleted"
// @Failure 500 {object} map[string]string "message: Error deleting tracks"
// @Router /track/range [delete]
// DeleteTrackByRange deletes tracks within a play count range
// from the user's library
func DeleteTrackByRange(c *gin.Context) {
	spotifyClient := c.MustGet("spotifyClient").(*spotifyAPI.Client)

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

	// delete all tracks from artists present
	// in the summary object
	for artistObj := range artistSummaryFiltered {
		tracksFilters, errTracks := userService.GetAllUserTracksByArtist(spotifyAPI.ID(rune(artistObj)), userTracks)

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

// GetTrackSummary godoc
// @Summary Get artists track count summary
// @Schemes
// @Description Returns a summary of tracks per artist, sorted by track count
// @Tags track
// @Accept json
// @Produce json
// @Param min query integer false "Minimum track count filter"
// @Param max query integer false "Maximum track count filter"
// @Success 200 {array} trackHelper.ArtistSummary
// @Failure 500 {object} map[string]string "message: Error getting tracks"
// @Router /track/summary [get]
func GetTrackSummary(c *gin.Context) {
	spotifyClient := c.MustGet("spotifyClient").(*spotifyAPI.Client)

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
