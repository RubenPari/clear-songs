package services

import (
	"errors"

	playlisthelper "github.com/RubenPari/clear-songs/src/helpers"
	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

func GetAllPlaylistTracks(id spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	if playlisthelper.CheckIfValidId(id) {
		return nil, errors.New("invalid playlist ID")
	}

	playlist, errGetPlaylist := utils.SpotifyClient.GetPlaylist(id)

	if errGetPlaylist != nil {
		return nil, errGetPlaylist
	}

	// get all tracks from playlist with pagination
	var offset = 0
	limit := 100
	var playlistTracks []spotifyAPI.PlaylistTrack

	for {
		tracks, errGetTracks := utils.SpotifyClient.GetPlaylistTracksOpt(playlist.ID, &spotifyAPI.Options{
			Offset: &offset,
			Limit:  &limit,
		}, "")

		if errGetTracks != nil {
			return nil, errGetTracks
		}

		playlistTracks = append(playlistTracks, tracks.Tracks...)

		if len(tracks.Tracks) < limit {
			break
		}

		offset += limit
	}

	return playlistTracks, nil
}

func DeleteTracksPlaylist(id spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack) error {
	if playlisthelper.CheckIfValidId(id) {
		return errors.New("invalid playlist ID")
	}

	trackIDs, errConvertIDs := utils.ConvertTracksToID(tracks)
	if errConvertIDs != nil {
		return errConvertIDs
	}

	// remove tracks from playlist 100 at a time
	for i := 0; i < len(trackIDs); i += 100 {
		end := i + 100

		if end > len(trackIDs) {
			end = len(trackIDs)
		}

		tracks100 := trackIDs[i:end]

		_, errDeleteTracks := utils.SpotifyClient.RemoveTracksFromPlaylist(id, tracks100...)

		if errDeleteTracks != nil {
			return errDeleteTracks
		}
	}

	return nil
}
