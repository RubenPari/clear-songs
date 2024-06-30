package services

import (
	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

func GetAllPlaylistTracks(idPlaylist spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	playlist, errGetPlaylist := utils.SpotifyClient.GetPlaylist(idPlaylist)

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

func DeleteTracksPlaylist(idPlaylist spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack) error {
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

		_, errDeleteTracks := utils.SpotifyClient.RemoveTracksFromPlaylist(idPlaylist, tracks100...)

		if errDeleteTracks != nil {
			return errDeleteTracks
		}
	}

	return nil
}
