package playlist

import (
	"log"
	"os"

	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/RubenPari/clear-songs/src/modules/utils"
	"github.com/gofiber/fiber/v2"
	spotifyAPI "github.com/zmb3/spotify/v2"
	"golang.org/x/exp/slices"
)

// CreateRapPlaylist is a function that creates
// a playlist with all rap songs in my library
func CreateRapPlaylist(c *fiber.Ctx) error {
	spotifyClient := authMO.SpotifyClient
	ctx := c.Context()

	// get all songs in my library
	allTracks, errAllTracks := utils.GetTracksUser()

	if errAllTracks != nil {
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't get all tracks",
		})
	}

	// filter rap songs
	var rapTracks []spotifyAPI.ID

	for _, track := range allTracks {
		artist, _ := spotifyClient.GetArtist(ctx, track.Artists[0].ID)

		if slices.Contains(artist.Genres, "rap") ||
			slices.Contains(artist.Genres, "hip hop") {
			rapTracks = append(rapTracks, track.ID)
		}
	}

	// get all songs in the playlist
	playlist, _ := spotifyClient.GetPlaylist(ctx, spotifyAPI.ID(os.Getenv("PLAYLIST_RAP")))

	// create an array with all songs id to remove in the playlist
	songsToRemove := make([]spotifyAPI.ID, 0)

	playlistTracks, _ := spotifyClient.GetPlaylistItems(ctx, playlist.ID)

	for _, playlistItem := range playlistTracks.Items {
		songsToRemove = append(songsToRemove, playlistItem.Track.Track.ID)
	}

	// remove 50 songs at time
	errRemoveTracks := utils.RemoveUserTracks(songsToRemove)

	if errRemoveTracks != nil {
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't remove songs",
		})
	}

	// add 100 songs at time to the playlist
	for start := 0; start < len(rapTracks); start += 100 {
		end := start + 50

		if end > len(rapTracks) {
			end = len(rapTracks)
		}

		_, errAddSongs := spotifyClient.AddTracksToPlaylist(ctx, playlist.ID, rapTracks[start:end]...)

		if errAddSongs != nil {
			log.Default().Printf("couldn't add songs: %v", errAddSongs)
			_ = c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"status":  "error",
				"message": "couldn't add songs",
			})
		}
	}

	_ = c.SendStatus(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "playlist created",
	})
}

// CreateEdmPlaylist is a function that creates
// a playlist with all edm songs in my library
func CreateEdmPlaylist(c *fiber.Ctx) error {
	return c.SendString("Create edm playlist")
}
