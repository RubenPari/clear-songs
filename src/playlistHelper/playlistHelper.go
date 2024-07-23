package playlisthelper

import (
	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

func CheckIfValidId(id spotifyAPI.ID) bool {
	if id.String() == "" {
		return false
	}

	_, errGetPlaylist := utils.SpotifyClient.GetPlaylist(id)

	return errGetPlaylist == nil
}
