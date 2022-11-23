package utils

import (
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify/v2"
)

// GetFeaturing get string containing all
// name of the artist featuring separated by ', '
// @param track spotifyAPI.SavedTrack
func GetFeaturing(track spotifyAPI.SavedTrack) string {
	var featuring = ""

	for _, artist := range track.Artists[1:] {
		featuring += artist.Name + ", "
	}

	return featuring
}

// ArrayArtistFoundedContains check if an array contains a value
func ArrayArtistFoundedContains(array *[100]string, value string) bool {
	for _, v := range *array {
		if v == value {
			return true
		}
	}

	return false
}

// AppendArray insert a value in
// the first empty position of an array
func AppendArray(array *[100]string, value string) {
	for i, v := range *array {
		if v == "" {
			(*array)[i] = value
			break
		}
	}
}

// ArrayArtistSummaryContains check if an array
// of ArtistLibrarySummary contains a value
func ArrayArtistSummaryContains(array *[]models.ArtistLibrarySummary, value string) (bool, int) {
	for i, v := range *array {
		if v.Id == value {
			return true, i
		}
	}

	return false, 0
}
