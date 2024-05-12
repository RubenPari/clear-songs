package services

import (
	"log"

	"github.com/RubenPari/clear-songs/src/utils"
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
// by artist id
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
// specified tracks from user
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

// GetAllUserAlbums
// returns all albums of user
func GetAllUserAlbums() []spotifyAPI.SavedAlbum {
	log.Default().Println("Getting all user albums")

	var allAlbums []spotifyAPI.SavedAlbum
	var offset = 0
	var limit = 50

	for {
		albums, err := utils.SpotifyClient.CurrentUsersAlbumsOpt(&spotifyAPI.Options{
			Limit:  &limit,
			Offset: &offset,
		})

		log.Default().Println("Getting albums from offset: ", offset)

		if err != nil {
			log.Default().Println("Error getting user albums")
			return nil
		}

		if len(albums.Albums) == 0 {
			break
		}

		allAlbums = append(allAlbums, albums.Albums...)

		offset += 50
	}

	log.Println("Total albums: ", len(allAlbums))

	return allAlbums
}

// GetAllUserAlbumsByArtist
// returns all albums of user
// by artist id
func GetAllUserAlbumsByArtist(idArtist spotifyAPI.ID) []spotifyAPI.SavedAlbum {
	log.Default().Println("Getting all user albums by artist")

	var filteredAlbums []spotifyAPI.SavedAlbum
	var offset = 0
	var limit = 50

	for {
		albums, err := utils.SpotifyClient.CurrentUsersAlbumsOpt(&spotifyAPI.Options{
			Limit:  &limit,
			Offset: &offset,
		})

		log.Default().Println("Getting albums from offset: ", offset)

		if err != nil {
			log.Default().Println("Error getting user albums")
			return nil
		}

		if len(albums.Albums) == 0 {
			break
		}

		// filter by artist id
		for _, album := range albums.Albums {
			if album.Artists[0].ID == idArtist {
				filteredAlbums = append(filteredAlbums, album)
				log.Default().Println("Album: ", album.Name, " - ", album.Artists[0].Name, " founded")
			}
		}

		offset += 50
	}

	log.Println("Total albums: ", len(filteredAlbums))

	return filteredAlbums
}

// ConvertAlbumToSongs
// converts album to songs
func ConvertAlbumToSongs(idAlbum spotifyAPI.ID) error {
	log.Default().Println("Converting album to songs")

	// get album info
	album, errAlbum := utils.SpotifyClient.GetAlbum(idAlbum)

	if errAlbum != nil {
		return errAlbum
	}

	// get all tracks of album
	var tracks []spotifyAPI.ID

	for _, track := range album.Tracks.Tracks {
		tracks = append(tracks, track.ID)
	}

	// add tracks to user library
	errAddTracks := utils.SpotifyClient.AddTracksToLibrary(tracks...)

	if errAddTracks != nil {
		return errAddTracks
	}

	return nil
}
