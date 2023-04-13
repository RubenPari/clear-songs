package artist

import (
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetArtistsSummary returns a
// map with the number of tracks
// of each artist
func GetArtistsSummary(tracks []spotifyAPI.SavedTrack) []models.ArtistSummary {
	var artistSummary = make(map[string]int)

	for _, track := range tracks {
		artistSummary[track.Artists[0].Name]++
	}

	var artistSummaryArray []models.ArtistSummary

	for artist, count := range artistSummary {
		artistSummaryArray = append(artistSummaryArray, models.ArtistSummary{
			Name:  artist,
			Count: count,
		})
	}

	return artistSummaryArray
}
