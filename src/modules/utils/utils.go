package utils

import (
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

// ArrayContains check if an array contains a value
func ArrayContains(array [50]string, value string) (bool, int) {
	for i, v := range array {
		if v == value {
			return true, i
		}
	}

	return false, 0
}

// AppendArray insert a value in
// the first empty position of an array
func AppendArray(array *[50]string, value string) {
	for i, v := range *array {
		if v == "" {
			(*array)[i] = value
			return
		}
	}
}
