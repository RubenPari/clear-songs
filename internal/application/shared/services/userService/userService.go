package userService

import (
	"errors"
	"log"

	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetAllUserTracks retrieves all tracks saved by the user.
//
// The function retrieves all tracks saved by the user, with pagination.
// The number of tracks retrieved at each iteration is limited to 50.
// The function returns a slice of spotifyAPI.SavedTrack and an error if the
// operation fails.
func GetAllUserTracks() ([]spotifyAPI.SavedTrack, error) {
	log.Default().Println("Getting all user tracks")

	var allTracks []spotifyAPI.SavedTrack
	var offset = 0
	var limit = 50

	for {
		tracks, err := utils.SpotifySvc.GetSpotifyClient().CurrentUsersTracksOpt(&spotifyAPI.Options{
			Limit:  &limit,
			Offset: &offset,
		})

		log.Default().Println("Getting tracks from offset: ", offset)

		if err != nil {
			log.Default().Println("Error getting user tracks")
			return nil, err
		}

		if len(tracks.Tracks) == 0 {
			break
		}

		allTracks = append(allTracks, tracks.Tracks...)

		offset += 50
	}

	log.Println("Total tracks: ", len(allTracks))

	return allTracks, nil
}

// GetAllUserTracksByArtist retrieves all tracks from a given artist saved by the user.
//
// id is the unique identifier of the Spotify artist.
// tracks is a slice of spotifyAPI.SavedTrack.
// The function returns a slice of spotifyAPI.ID and an error if the
// operation fails.
func GetAllUserTracksByArtist(id spotifyAPI.ID, tracks []spotifyAPI.SavedTrack) ([]spotifyAPI.ID, error) {
	log.Default().Println("Getting all user tracks by artist")

	var filteredTracks []spotifyAPI.ID

	for _, track := range tracks {
		if track.Artists[0].ID == id {
			filteredTracks = append(filteredTracks, track.ID)
			log.Default().Println("Track: ", track.Name, " - ", track.Artists[0].Name, " founded")
		}
	}

	log.Println("Total tracks: ", len(filteredTracks))

	return filteredTracks, nil
}

// GetUserTracksByArtist retrieves all full track objects from a given artist saved by the user.
//
// id is the unique identifier of the Spotify artist.
// tracks is a slice of spotifyAPI.SavedTrack.
// The function returns a slice of spotifyAPI.SavedTrack and an error if the operation fails.
func GetUserTracksByArtist(id spotifyAPI.ID, tracks []spotifyAPI.SavedTrack) ([]spotifyAPI.SavedTrack, error) {
	log.Default().Println("Getting all user track objects by artist")

	var filteredTracks []spotifyAPI.SavedTrack

	for _, track := range tracks {
		if track.Artists[0].ID == id {
			filteredTracks = append(filteredTracks, track)
		}
	}

	log.Println("Total tracks found: ", len(filteredTracks))

	return filteredTracks, nil
}

// DeleteTracksUser removes a list of tracks from the user's library.
//
// tracks is a slice of spotifyAPI.ID representing the track IDs to be deleted.
// The function processes the deletions in batches of 50 tracks at a time.
// If an error occurs while deleting the tracks, it returns the error.
func DeleteTracksUser(c *gin.Context, tracks []spotifyAPI.ID) error {
	service := utils.GetSpotifyService(c)
	if service == nil {
		return errors.New("spotify service not available")
	}

	client := service.GetSpotifyClient()
	if client == nil {
		return errors.New("spotify client not available")
	}

	log.Default().Println("Deleting user tracks")

	var offset = 0
	var limit = 50

	for {
		if offset >= len(tracks) {
			break
		}

		if offset+50 > len(tracks) {
			limit = len(tracks) - offset
		}

		err := client.RemoveTracksFromLibrary(tracks[offset : offset+limit]...)

		if err != nil {
			log.Default().Println("Error deleting user tracks")
			return err
		}

		log.Default().Println("Deleting tracks from offset: ", offset)

		offset += 50
	}

	return nil
}
