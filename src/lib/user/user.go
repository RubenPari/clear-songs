package user

import (
	"log"

	"github.com/RubenPari/clear-songs/src/lib/array"
	"github.com/RubenPari/clear-songs/src/lib/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetAllUserTracks
// returns all tracks
// of user
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

// GetAllUserTracksByArtist
// returns all tracks of user
func GetAllUserTracksByArtist(id spotifyAPI.ID) ([]spotifyAPI.ID, error) {
	log.Default().Println("Getting all user tracks by artist")

	var filteredTracks []spotifyAPI.ID
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

		// filter by artist id
		for _, track := range tracks.Tracks {
			if track.Artists[0].ID == id {
				filteredTracks = append(filteredTracks, track.ID)
				log.Default().Println("Track: ", track.Name, " - ", track.Artists[0].Name, " founded")
			}
		}

		offset += 50
	}

	log.Println("Total tracks: ", len(filteredTracks))

	return filteredTracks, nil
}

// DeleteTracksUser deletes
// all tracks of user
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

// GetAllUserTracksByGenre return all
// user track library by genre
func GetAllUserTracksByGenre(genre string) ([]spotifyAPI.ID, error) {
	// get all possible genres name
	genres := utils.GetPossibleGenres(genre)

	var tracksFilter []spotifyAPI.ID

	var offset = 0
	var limit = 50

	log.Default().Println("Getting all user tracks by genre")

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

		// filter by genre name
		for _, track := range tracks.Tracks {
			// get artist info object
			artist, _ := utils.SpotifyClient.GetArtist(track.Artists[0].ID)

			// check if artist has the specific genre
			if array.ContainsGenre(artist.Genres, genres) {
				tracksFilter = append(tracksFilter, track.ID)
			}
		}

		offset += 50
	}

	return tracksFilter, nil
}

func DeleteTracksByArtists(artists []spotifyAPI.FullArtist) error {
	log.Default().Println("Deleting tracks by artists")

	var tracks []spotifyAPI.ID

	// get all tracks of user
	allTracks, err := GetAllUserTracks()

	if err != nil {
		return err
	}

	// filter tracks by artist
	for _, track := range allTracks {
		for _, artist := range artists {
			if track.Artists[0].ID == artist.ID {
				tracks = append(tracks, track.ID)
			}
		}
	}

	// delete tracks
	err = DeleteTracksUser(tracks)

	if err != nil {
		return err
	}

	return nil
}
