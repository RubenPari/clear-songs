package artist

import spotifyAPI "github.com/zmb3/spotify"

// GetArtistsSummary returns a
// map with the number of tracks
// of each artist
func GetArtistsSummary(tracks []spotifyAPI.SavedTrack) map[string]int {
	var artistSummary = make(map[string]int)

	for _, track := range tracks {
		artistSummary[track.Artists[0].Name]++
	}

	return artistSummary
}
