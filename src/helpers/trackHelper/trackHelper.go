package trackHelper

import (
	"log"

	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetArtistsSummary generates an array of ArtistSummary based on the tracks provided.
// It counts the number of tracks for each artist and creates a summary with the artist's name, ID, and track count.
//
// Parameters:
//   - tracks: a slice of spotifyAPI.SavedTrack representing the tracks to analyze
//
// Returns:
//   - []models.ArtistSummary: an array of ArtistSummary with name, ID, and track count for each artist
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
