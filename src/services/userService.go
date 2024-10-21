package services

import (
	"log"

	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetAllUserTracks retrieves all tracks saved by the user.
//
// The function retrieves all tracks saved by the user, with pagination.
// The number of tracks retrieved at each iteration is limited to 50.
// The function returns a slice of spotifyAPI.SavedTrack and an error if the
// operation fails.
func GetAllUserTracks() ([]spotifyAPI.SavedTrack, error) {
	var allTracks []spotifyAPI.SavedTrack
	var offset = 0
	var limit = 50

	log.Default().Println("Getting all user tracks")

	for {
		tracks, err := utils.SpotifyClient.CurrentUsersTracksOpt(&spotifyAPI.Options{
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

// DeleteTracksUser removes a list of tracks from the user's library.
//
// tracks is a slice of spotifyAPI.ID representing the track IDs to be deleted.
// The function processes the deletions in batches of 50 tracks at a time.
// If an error occurs while deleting the tracks, it returns the error.
func DeleteTracksUser(tracks []spotifyAPI.ID) error {
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

		err := utils.SpotifyClient.RemoveTracksFromLibrary(tracks[offset : offset+limit]...)

		if err != nil {
			log.Default().Println("Error deleting user tracks")
			return err
		}

		log.Default().Println("Deleting tracks from offset: ", offset)

		offset += 50
	}

	return nil
}

// ConvertAlbumToSongs converts an album to songs and saves them to the user library.
//
// idAlbum is the unique identifier of the Spotify album.
// The function first retrieves the album's info and then all its tracks.
// After that, it adds the tracks to the user library.
// If an error occurs while converting the album, it returns the error.
func ConvertAlbumToSongs(idAlbum spotifyAPI.ID) error {
	log.Default().Println("Converting album to songs")

	// get album info
	album, errAlbum := utils.SpotifyClient.GetAlbum(idAlbum)

	if errAlbum != nil {
		return errAlbum
	}

	// get all tracks of album
	trackIDs, errConvertIDs := utils.ConvertTracksToID(album.Tracks.Tracks)
	if errConvertIDs != nil {
		return errConvertIDs
	}

	// add tracks to user library
	errAddTracks := utils.SpotifyClient.AddTracksToLibrary(trackIDs...)

	if errAddTracks != nil {
		return errAddTracks
	}

	return nil
}
