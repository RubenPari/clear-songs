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
