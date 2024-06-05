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

	firstTracks, errGetTracks := utils.SpotifyClient.GetPlaylistTracks(playlist.ID)

	if errGetTracks != nil {
		return nil, errGetTracks
	}

	// get all tracks from playlist with pagination
	offset := 0
	limit := 100

	for {
		tracks, errGetTracks := utils.SpotifyClient.GetPlaylistTracksOpt(playlist.ID, &spotifyAPI.Options{
			Offset: &offset,
			Limit:  &limit,
		}, "")

		if errGetTracks != nil {
			return nil, errGetTracks
		}

		if len(tracks.Tracks) == 0 {
			break
		}

		firstTracks.Tracks = append(firstTracks.Tracks, tracks.Tracks...)

		offset += 100
	}
	return firstTracks.Tracks, nil
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
