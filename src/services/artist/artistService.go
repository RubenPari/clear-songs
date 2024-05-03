package artist

import (
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
	"log"
)

// GetArtistsSummary returns a
// map with the number of tracks
// of each artist
func GetArtistsSummary(tracks []spotifyAPI.SavedTrack) []models.ArtistSummary {
	log.Default().Println("Getting artists summary array")

	var artistSummary = make(map[string]struct {
		count int
		id    string
	})

	for _, track := range tracks {
		// Check if artist is already in the map
		if artist, exists := artistSummary[track.Artists[0].Name]; exists {
			artist.count++
			artistSummary[track.Artists[0].Name] = artist
		} else {
			artistSummary[track.Artists[0].Name] = struct {
				count int
				id    string
			}{
				count: 1,
				id:    string(track.Artists[0].ID),
			}
		}
	}

	var artistSummaryArray []models.ArtistSummary

	for artist, summary := range artistSummary {
		artistSummaryArray = append(artistSummaryArray, models.ArtistSummary{
			Name:  artist,
			Id:    summary.id,
			Count: summary.count,
		})
	}

	return artistSummaryArray
}
