package artist

import (
	"io/ioutil"
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
	// get file content
	file, errOpen := FileHeader.Open()

	if errOpen != nil {
		return nil, errOpen
	}

	defer file.Close()

	content, errRead := ioutil.ReadAll(file)

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
		artistobj, errSearch := utils.SpotifyClient.Search(artist, spotifyAPI.SearchTypeArtist)

		if errSearch != nil {
			return nil, errSearch
		}

		// append the artist to the slice
		artists = append(artists, artistobj.Artists.Artists[0])
	}

	return artists, nil
}
