package artist

import (
	"io"
	"log"
	"mime/multipart"
	"strings"

	"github.com/RubenPari/clear-songs/src/models"
	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
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
		// Controllo se l'artista è già nella mappa
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

// GroupArtistSummaryByGenres returns a
// map of []models.ArtistSummary grouped
// by genres musical
func GroupArtistSummaryByGenres(artistSummaryArray []models.ArtistSummary) []models.ArtistGroupSummary {
	log.Default().Println("Grouping artists summary by genres")

	var artistSummaryGrouped = make(map[string][]models.ArtistSummary)

	for _, artistSummary := range artistSummaryArray {
		artist := artistSummary.Name

		// get artist spotify object by searching its name
		artistObj, errSearch := utils.SpotifyClient.Search(artist, spotifyAPI.SearchTypeArtist)

		if errSearch != nil {
			continue
		}

		// get artist genres
		artistGenres := artistObj.Artists.Artists[0].Genres

		// add artist summary to the map
		for _, genre := range artistGenres {
			artistSummaryGrouped[genre] = append(artistSummaryGrouped[genre], artistSummary)
		}
	}

	var artistSummaryGroupedArray []models.ArtistGroupSummary

	for genre, artistSummary := range artistSummaryGrouped {
		artistSummaryGroupedArray = append(artistSummaryGroupedArray, models.ArtistGroupSummary{
			Genre:   genre,
			Artists: artistSummary,
		})
	}

	return artistSummaryGroupedArray
}

// GetArtistsFromFile returns an array
// of spotifyAPI.FullArtist from a .txt
// file that contains the name of the
// artists separated by new line
func GetArtistsFromFile(FileHeader *multipart.FileHeader) ([]spotifyAPI.FullArtist, error) {
	log.Default().Println("Getting artists from file")

	// get file content
	file, errOpen := FileHeader.Open()

	if errOpen != nil {
		return nil, errOpen
	}

	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	content, errRead := io.ReadAll(file)

	if errRead != nil {
		return nil, errRead
	}

	// convert the content of []byte to string
	contentString := string(content)

	// remove the last new line
	contentString = strings.TrimSuffix(contentString, "\n")

	// cut the string by new line
	artistsFile := strings.Split(contentString, "\n")

	// create a slice of spotifyAPI.SimpleArtist
	artists := make([]spotifyAPI.FullArtist, 0)

	for _, artist := range artistsFile {
		// get artist spotify object by searching its name
		artistObj, errSearch := utils.SpotifyClient.Search(artist, spotifyAPI.SearchTypeArtist)

		if errSearch != nil {
			return nil, errSearch
		}

		// append the artist to the slice
		artists = append(artists, artistObj.Artists.Artists[0])
	}

	return artists, nil
}
