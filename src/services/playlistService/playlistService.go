package playlistService

import (
	"errors"

	"github.com/RubenPari/clear-songs/src/helpers/playlistHelper"

	"github.com/RubenPari/clear-songs/src/constants"

	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetAllPlaylistTracks retrieves all tracks from a Spotify playlist.
//
// id is the unique identifier of the Spotify playlist.
// Returns a slice of spotifyAPI.PlaylistTrack and an error if the operation fails.
func GetAllPlaylistTracks(id spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	if !playlistHelper.CheckIfValidId(id) {
		return nil, errors.New("invalid playlist ID")
	}

	playlist, errGetPlaylist := utils.SpotifySvc.GetSpotifyClient().GetPlaylist(id)

	if errGetPlaylist != nil {
		return nil, errGetPlaylist
	}

	// get all tracks from playlist with pagination
	limit := constants.LimitGetPlaylistTracks
	offset := constants.Offset
	var playlistTracks []spotifyAPI.PlaylistTrack

	for {
		tracks, errGetTracks := utils.SpotifySvc.GetSpotifyClient().GetPlaylistTracksOpt(playlist.ID, &spotifyAPI.Options{
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

// GetAllUserPlaylists retrieves all playlists owned or followed by the user.
//
// Returns a slice of spotifyAPI.SimplePlaylist and an error if the operation fails.
func GetAllUserPlaylists() ([]spotifyAPI.SimplePlaylist, error) {
	var allPlaylists []spotifyAPI.SimplePlaylist
	var offset = 0
	var limit = 50

	for {
		playlists, err := utils.SpotifySvc.GetSpotifyClient().CurrentUsersPlaylistsOpt(&spotifyAPI.Options{
			Limit:  &limit,
			Offset: &offset,
		})

		if err != nil {
			return nil, err
		}

		if len(playlists.Playlists) == 0 {
			break
		}

		allPlaylists = append(allPlaylists, playlists.Playlists...)

		if len(playlists.Playlists) < limit {
			break
		}

		offset += limit
	}

	return allPlaylists, nil
}

// DeletePlaylistTracks deletes tracks from a Spotify playlist.
//
// id is the unique identifier of the Spotify playlist and tracks is a slice of spotifyAPI.PlaylistTrack to be deleted.
// Returns an error if the operation fails.
func DeletePlaylistTracks(id spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack) error {
	if !playlistHelper.CheckIfValidId(id) {
		return errors.New("invalid playlist ID")
	}

	trackIDs, errConvertIDs := utils.ConvertTracksToID(tracks)
	if errConvertIDs != nil {
		return errConvertIDs
	}

	// remove tracks from playlist
	limit := constants.LimitRemovePlaylistTracks
	offset := constants.Offset

	for i := offset; i < len(trackIDs); i += limit {
		end := i + limit

		if end > len(trackIDs) {
			end = len(trackIDs)
		}

		tracksPagination := trackIDs[i:end]

		_, errDeleteTracks := utils.SpotifySvc.GetSpotifyClient().RemoveTracksFromPlaylist(id, tracksPagination...)

		if errDeleteTracks != nil {
			return errDeleteTracks
		}
	}

	return nil
}
