package artist

import (
	"io"
	"log"
	"mime/multipart"
	"strings"

	"github.com/RubenPari/clear-songs/src/lib/utils"
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetArtistsSummary returns a
// map with the number of tracks
// of each artist
func GetArtistsSummary(tracks []spotifyAPI.SavedTrack) []models.ArtistSummary {
	log.Default().Println("Getting artists summary array")

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
