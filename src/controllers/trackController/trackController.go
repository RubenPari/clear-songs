/**
 * Track Controller Package
 * 
 * This package handles all track management HTTP endpoints. It provides
 * functionality for retrieving track summaries, deleting tracks by artist,
 * and deleting tracks by count range.
 * 
 * All endpoints require Spotify authentication via SpotifyAuthMiddleware,
 * which ensures the user has a valid OAuth token before processing requests.
 * 
 * Features:
 * - Get track summaries grouped by artist with optional filtering
 * - Delete all tracks from a specific artist (with backup)
 * - Delete tracks based on count ranges (e.g., artists with 1-5 tracks)
 * 
 * The controller uses caching to improve performance and reduce Spotify API calls.
 * Track data is cached and invalidated when modifications are made.
 * 
 * @package trackController
 * @author Clear Songs Development Team
 */
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

/**
 * DeleteTrackByArtist deletes all tracks from a specific artist
 * 
 * This endpoint removes all tracks from the user's Spotify library that belong
 * to the specified artist. The operation:
 * 1. Retrieves cached user tracks (or fetches from Spotify if not cached)
 * 2. Filters tracks by the specified artist ID
 * 3. Creates a backup of tracks to be deleted
 * 4. Removes tracks from user's Spotify library
 * 5. Updates the database and cache
 * 
 * The artist ID is provided as a URL parameter: /track/by-artist/:id_artist
 * 
 * @param c - Gin context containing HTTP request and response
 */
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

	// Get Spotify client from context (set by SpotifyAuthMiddleware)
	var spotifyClient *spotifyAPI.Client
	if client, exists := c.Get("spotifyClient"); exists {
		spotifyClient = client.(*spotifyAPI.Client)
	}

	artistSummaryArray := trackHelper.GetArtistsSummary(userTracks, spotifyClient)
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

	// Get Spotify client from context (set by SpotifyAuthMiddleware)
	var spotifyClient *spotifyAPI.Client
	if client, exists := c.Get("spotifyClient"); exists {
		spotifyClient = client.(*spotifyAPI.Client)
	}

	artistSummaryArray := trackHelper.GetArtistsSummary(tracks, spotifyClient)

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
