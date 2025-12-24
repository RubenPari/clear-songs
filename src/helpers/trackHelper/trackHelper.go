package trackHelper

import (
	"log"

	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetArtistsSummary generates an array of ArtistSummary based on the tracks provided.
// It counts the number of tracks for each artist and creates a summary with the artist's name, ID, track count, and image URL.
//
// Parameters:
//   - tracks: a slice of spotifyAPI.SavedTrack representing the tracks to analyze
//   - client: optional Spotify client for fetching artist images (can be nil)
//
// Returns:
//   - []models.ArtistSummary: an array of ArtistSummary with name, ID, track count, and image URL for each artist
func GetArtistsSummary(tracks []spotifyAPI.SavedTrack, client *spotifyAPI.Client) []models.ArtistSummary {
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
		imageURL := ""
		
		// Fetch artist image if client is available
		if client != nil && summary.id != "" {
			artistInfo, err := client.GetArtist(spotifyAPI.ID(summary.id))
			if err == nil && len(artistInfo.Images) > 0 {
				// Use the smallest image (usually the last one) for better performance
				// or medium image if available
				for i := len(artistInfo.Images) - 1; i >= 0; i-- {
					if artistInfo.Images[i].Width <= 300 || i == 0 {
						imageURL = artistInfo.Images[i].URL
						break
					}
				}
			}
		}

		artistSummaryArray = append(artistSummaryArray, models.ArtistSummary{
			Name:     artist,
			Id:       summary.id,
			Count:    summary.count,
			ImageURL: imageURL,
		})
	}

	return artistSummaryArray
}
