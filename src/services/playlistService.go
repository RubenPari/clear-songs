package services

import (
	"errors"
	"log"

	"github.com/RubenPari/clear-songs/src/constants"

	"github.com/RubenPari/clear-songs/src/helpers"
	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

// GetAllPlaylistTracks retrieves all tracks from a Spotify playlist.
//
// id is the unique identifier of the Spotify playlist.
// Returns a slice of spotifyAPI.PlaylistTrack and an error if the operation fails.
func GetAllPlaylistTracks(id spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	if !helpers.CheckIfValidId(id) {
		return nil, errors.New("invalid playlist ID")
	}

	playlist, errGetPlaylist := utils.SpotifyClient.GetPlaylist(id)

	if errGetPlaylist != nil {
		return nil, errGetPlaylist
	}

	// get all tracks from playlist with pagination
	limit := constants.LimitGetPlaylistTracks
	offset := constants.Offset
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

// DeletePlaylistTracks deletes tracks from a Spotify playlist.
//
// id is the unique identifier of the Spotify playlist and tracks is a slice of spotifyAPI.PlaylistTrack to be deleted.
// Returns an error if the operation fails.
func DeletePlaylistTracks(id spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack) error {
	if !helpers.CheckIfValidId(id) {
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

		_, errDeleteTracks := utils.SpotifyClient.RemoveTracksFromPlaylist(id, tracksPagination...)

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
	log.Default().Println("Creating playlist with minor songs")

	userId, _ := utils.GetUserId()
	var idPlaylistMinorSongs *spotifyAPI.ID

	// check if "MinorSongs" playlist already exists
	playlistsUser, errPlaylistsUser := utils.SpotifyClient.GetPlaylistsForUser(string(userId))

	if errPlaylistsUser != nil {
		return errPlaylistsUser
	}

	for _, playlist := range playlistsUser.Playlists {
		if playlist.Name == constants.PlaylistNameWithMinorSongs {
			idPlaylistMinorSongs = &playlist.ID
		}
	}

	if idPlaylistMinorSongs == nil {
		// create playlist
		playlistMinorSongs, errCreate := utils.SpotifyClient.CreatePlaylistForUser(
			string(userId),
			constants.PlaylistNameWithMinorSongs,
			constants.DescriptionPlaylistNameWithMinorSongs,
			false,
		)

		if errCreate != nil {
			return errCreate
		}

		idPlaylistMinorSongs = &playlistMinorSongs.ID
	}

	artistsSummary := helpers.GetArtistsSummary(tracks)

	// find all tracks that belong to artists with less than 5 songs
	var tracksToKeep []spotifyAPI.SavedTrack
	for _, artistSummary := range artistsSummary {
		if artistSummary.Count <= 5 {
			log.Default().Printf("Artist %s has less than 5 songs, start getting tracks", artistSummary.Name)
			for _, track := range tracks {
				if track.Artists[0].ID == spotifyAPI.ID(artistSummary.Id) {
					tracksToKeep = append(tracksToKeep, track)
				}
			}
			log.Default().Printf("Getted %d tracks from artist %s", len(tracksToKeep), artistSummary.Name)
		}
	}

	// convert tracks to IDs
	trackIDs, errConvertIDs := utils.ConvertTracksToID(tracksToKeep)
	if errConvertIDs != nil {
		return errConvertIDs
	}

	// insert tracks into playlist with pagination
	var offset = constants.Offset
	limit := constants.LimitInsertPlaylistTracks

	for {
		_, errGetTracks := utils.SpotifyClient.AddTracksToPlaylist(*idPlaylistMinorSongs, trackIDs[offset:offset+limit]...)

		log.Default().Printf("Added tracks from offset %d to %d", offset, offset+limit)

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
