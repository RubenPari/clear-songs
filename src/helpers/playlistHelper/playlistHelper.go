package playlistHelper

import (
	"github.com/RubenPari/clear-songs/src/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

// CheckIfValidId checks if the provided Spotify playlist ID is valid.
// It takes the Spotify playlist ID as input and returns a boolean value indicating the validity.
//
// Parameters:
// - id: The unique identifier of the Spotify playlist (type: spotifyAPI.ID)
//
// Returns:
// - bool: true if the playlist ID is valid, false otherwise
func CheckIfValidId(id spotifyAPI.ID) bool {
	if id.String() == "" {
		return false
	}

	_, errGetPlaylist := utils.SpotifySvc.GetSpotifyClient().GetPlaylist(id)

	return errGetPlaylist == nil
}
