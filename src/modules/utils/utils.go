package utils

import (
	"context"
	"log"

	"github.com/RubenPari/clear-songs/src/models"
	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
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

// ArrayArtistFoundedContains check if an array contains a value
func ArrayArtistFoundedContains(array *[100]string, value string) bool {
	for _, v := range *array {
		if v == value {
			return true
		}
	}

	return false
}

// AppendArray insert a value in
// the first empty position of an array
func AppendArray(array *[100]string, value string) {
	for i, v := range *array {
		if v == "" {
			(*array)[i] = value
			break
		}
	}
}

// ArrayArtistSummaryContains check if an array
// of ArtistLibrarySummary contains a value
func ArrayArtistSummaryContains(array *[]models.ArtistLibrarySummary, value string) (bool, int) {
	for i, v := range *array {
		if v.Id == value {
			return true, i
		}
	}

	return false, 0
}

// GetTracksUser get all tracks saved by the user
func GetTracksUser() ([]spotifyAPI.SavedTrack, error) {
	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	// call to spotify api n times based on the offset
	var limit = 50
	var offset = 0 // NOTE: offset is excluded
	// NOTE: gli elementi selezionati da n a m range
	// non sono salvati in tale ordine

	var errGetSavedTracks error

	songs := make([]spotifyAPI.SavedTrack, 0)

	for {
		tracksPage, err := spotifyClient.CurrentUsersTracks(ctx, spotifyAPI.Limit(limit), spotifyAPI.Offset(offset))

		if err != nil {
			errGetSavedTracks = err
			break
		}

		songs = append(songs, tracksPage.Tracks...)

		if len(tracksPage.Tracks) < limit {
			break
		}

		offset += limit
	}

	if errGetSavedTracks != nil {
		log.Default().Printf("couldn't get songs: %v", errGetSavedTracks)
		return nil, errGetSavedTracks
	} else {
		return songs, nil
	}
}

func RemoveUserTracks(tracks []spotifyAPI.ID) error {
	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	var errRemoveTracks error = nil

	for start := 0; start < len(tracks); start += 50 {
		end := start + 50

		if end > len(tracks) {
			end = len(tracks)
		}

		err := spotifyClient.RemoveTracksFromLibrary(ctx, tracks[start:end]...)

		if errRemoveTracks != nil {
			log.Default().Printf("couldn't remove songs: %v", errRemoveTracks)
			errRemoveTracks = err
			break
		}
	}

	return errRemoveTracks
}
