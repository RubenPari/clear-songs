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

	// get all tracks from user
	var tracks []spotifyAPI.SavedTrack
	var errTracks error

	value, found := cacheManager.Get("userTracks")

	if found {
		tracks = value.([]spotifyAPI.SavedTrack)
	} else {
		tracks, errTracks = userService.GetAllUserTracks(spotifyClient)

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

	// get tracks from user
	var tracks []spotifyAPI.SavedTrack
	var errTracks error

	value, found := cacheManager.Get("userTracks")

	if found {
		tracks = value.([]spotifyAPI.SavedTrack)
	} else {
		tracks, errTracks = userService.GetAllUserTracks(spotifyClient)

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
