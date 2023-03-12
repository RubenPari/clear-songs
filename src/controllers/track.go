package controllers

import (
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
	"strconv"
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
	tracks, errTracks := utils.GetAllUserTracks()

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
	}

	// initialize artist summary array
	// with attributes: name, count
	artistSummaryArray := make(map[string]int)

	for _, page := range tracks {
		artistSummaryArray[page.Artists[0].Name]++
	}

	// filter by min
	if minStr != "" {
		artistSummaryArray = utils.FilterByMin(artistSummaryArray, min)
	}

	// filter by max
	if maxStr != "" {
		artistSummaryArray = utils.FilterByMax(artistSummaryArray, max)
	}

	c.JSON(200, artistSummaryArray)
}

// DeleteTrackByArtist deletes all tracks from an artist
func DeleteTrackByArtist(c *gin.Context) {
	// get spotify client
	spotify := utils.SpotifyClient

	// get artist id from url
	idArtist := c.Param("id_artist")

	// get tracks from user
	tracks, errTracks := spotify.CurrentUsersTracks()

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
	}

	// initialize tracks array
	tracksArrayToDelete := make([]spotifyAPI.ID, 0)

	for _, page := range tracks.Tracks {
		if page.Artists[0].ID.String() == idArtist {
			tracksArrayToDelete = append(tracksArrayToDelete, page.ID)
		}
	}

	// delete tracks from artist
	errDelete := spotify.RemoveTracksFromLibrary(tracksArrayToDelete...)

	if errDelete != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting tracks",
		})
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Tracks deleted",
	})
}

func DeleteTrackByGenre(c *gin.Context) {
	// get spotify client
	spotify := utils.SpotifyClient

	// get genre name from query params
	genre := c.Query("genre")

	// create genres array
	genres := utils.GetPossibleGenres(genre)

	// get tracks from user
	tracks, errTracks := spotify.CurrentUsersTracks()

	if errTracks != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting tracks",
		})
	}

	// initialize tracks array
	tracksArrayToDelete := make([]spotifyAPI.ID, 0)

	for _, page := range tracks.Tracks {
		// get artist genres from track's album
		fullAlbum, _ := spotify.GetAlbum(page.Album.ID)

		if utils.ContainsGenre(fullAlbum.Genres, genres) {
			tracksArrayToDelete = append(tracksArrayToDelete, page.ID)
		}
	}

	// delete tracks from artist
	errDelete := spotify.RemoveTracksFromLibrary(tracksArrayToDelete...)

	if errDelete != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting tracks",
		})
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Tracks deleted",
	})
}
