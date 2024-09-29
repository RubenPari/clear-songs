package services

import (
	"errors"

	playlisthelper "github.com/RubenPari/clear-songs/src/helpers"
	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetAllPlaylistTracks retrieves all tracks from a Spotify playlist.
//
// id is the unique identifier of the Spotify playlist.
// Returns a slice of spotifyAPI.PlaylistTrack and an error if the operation fails.
func GetAllPlaylistTracks(id spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	if !playlisthelper.CheckIfValidId(id) {
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

// DeleteTracksPlaylist deletes tracks from a Spotify playlist.
//
// id is the unique identifier of the Spotify playlist and tracks is a slice of spotifyAPI.PlaylistTrack to be deleted.
// Returns an error if the operation fails.
func DeleteTracksPlaylist(id spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack) error {
	if !playlisthelper.CheckIfValidId(id) {
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

// CreatePlaylistTracksMinor creates a playlist with all tracks from the user library
// that belong to artists for which the user owns a maximum of 5 tracks.
//
// tracks is a slice of spotifyAPI.SavedTrack to be filtered and added to the playlist.
// Returns an error if the operation fails.
func CreatePlaylistTracksMinor(tracks []spotifyAPI.SavedTrack) error {
	// get user id
	user, errorUser := utils.SpotifyClient.CurrentUser()

	if errorUser != nil {
		return errorUser
	}

	userId := user.ID

	var idPlaylistMinorSongs *spotifyAPI.ID

	// check if "MinorSongs" playlist already exists
	playlistsUser, errplaylistsUser := utils.SpotifyClient.GetPlaylistsForUser(userId)

	if errplaylistsUser != nil {
		return errplaylistsUser
	}

	for _, playlist := range playlistsUser.Playlists {
		if playlist.Name == "MinorSongs" {
			idPlaylistMinorSongs = &playlist.ID
		}
	}

	if idPlaylistMinorSongs == nil {
		// create playlist
		playlistMinorSongs, errCreate := utils.SpotifyClient.CreatePlaylistForUser(
			userId,
			"MinorSongs",
			"MinorSongs",
			false,
		)

		if errCreate != nil {
			return errCreate
		}

		idPlaylistMinorSongs = &playlistMinorSongs.ID
	}

	artistsSummary := utils.GetArtistsSummary(tracks)

	// find all tracks that belong to artists with less than 5 songs
	var tracksToKeep []spotifyAPI.SavedTrack
	for _, artistSummary := range artistsSummary {
		if artistSummary.Count <= 5 {
			for _, track := range tracks {
				if track.Artists[0].ID == spotifyAPI.ID(artistSummary.Id) {
					tracksToKeep = append(tracksToKeep, track)
				}
			}
		}
	}

	// convert tracks to IDs
	trackIDs, errConvertIDs := utils.ConvertTracksToID(tracksToKeep)
	if errConvertIDs != nil {
		return errConvertIDs
	}

	// insert tracks into playlist with pagination
	var offset = 0
	limit := 100

	for {
		_, errGetTracks := utils.SpotifyClient.AddTracksToPlaylist(*idPlaylistMinorSongs, trackIDs[offset:offset+limit]...)

		if errGetTracks != nil {
			return errGetTracks
		}

		if len(trackIDs[offset:offset+limit]) < limit {
			break
		}

		offset += limit
	}

	return nil
}
