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
	for firstTracks.Next != "" {
		offset := firstTracks.Offset + firstTracks.Limit
		nextTracks, errNextTracks := utils.SpotifyClient.GetPlaylistTracksOpt(playlist.ID, &spotifyAPI.Options{
			Offset: &offset,
		}, "")

		if errNextTracks != nil {
			return nil, errNextTracks
		}

		firstTracks.Tracks = append(firstTracks.Tracks, nextTracks.Tracks...)
		firstTracks.Next = nextTracks.Next
	}

	return firstTracks.Tracks, nil
}

func DeleteTracksPlaylist(idPlaylist spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack) error {
	// get all track ids
	var trackIDs []spotifyAPI.ID
	for _, track := range tracks {
		trackIDs = append(trackIDs, track.Track.ID)
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
