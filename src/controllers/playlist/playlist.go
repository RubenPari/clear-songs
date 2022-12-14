package playlist

import (
	"context"
	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/gofiber/fiber/v2"
	spotifyAPI "github.com/zmb3/spotify/v2"
	"golang.org/x/exp/slices"
	"log"
	"os"
)

// CreateRapPlaylist is a function that creates
// a playlist with all rap songs in my library
func CreateRapPlaylist(c *fiber.Ctx) error {
	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	// call to spotify api n times based on the offset
	var limit = 50
	var offset = 0 // NOTE: offset is excluded
	// NOTE: gli elementi selezionati da n a m range
	// non sono salvati in tale ordine

	var errGetSavedTracks error

	songsToAdd := make([]spotifyAPI.ID, 0)

	// get all songs in my library what have the genre rap
	for {
		tracksPage, err := spotifyClient.CurrentUsersTracks(ctx, spotifyAPI.Limit(limit), spotifyAPI.Offset(offset))

		if err != nil {
			errGetSavedTracks = err
			break
		}

		// for each group of songs check if the genre (of artist) is rap
		for _, track := range tracksPage.Tracks {
			artist, _ := spotifyClient.GetArtist(ctx, track.Artists[0].ID)

			if slices.Contains(artist.Genres, "rap") ||
				slices.Contains(artist.Genres, "hip hop") {
				songsToAdd = append(songsToAdd, track.ID)
			}
		}

		if len(tracksPage.Tracks) < limit {
			break
		}

		offset += limit
	}

	if errGetSavedTracks != nil {
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't get songs",
		})
	}

	// get all songs in the playlist
	playlist, _ := spotifyClient.GetPlaylist(ctx, spotifyAPI.ID(os.Getenv("PLAYLIST_RAP")))

	songsToRemove := make([]spotifyAPI.ID, 0)

	playlistTracks, _ := spotifyClient.GetPlaylistItems(ctx, playlist.ID)

	for _, playlistItem := range playlistTracks.Items {
		songsToRemove = append(songsToRemove, playlistItem.Track.Track.ID)
	}

	// remove 50 songs at time
	for start := 0; start < len(songsToRemove); start += 50 {
		end := start + 50

		if end > len(songsToRemove) {
			end = len(songsToRemove)
		}

		errRemoveSongs := spotifyClient.RemoveTracksFromLibrary(ctx, songsToRemove[start:end]...)

		if errRemoveSongs != nil {
			log.Default().Printf("couldn't remove songs: %v", errRemoveSongs)
			_ = c.SendStatus(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"status":  "error",
				"message": "couldn't remove songs",
			})
		}
	}

	// add 100 songs at time
	for start := 0; start < len(songsToAdd); start += 100 {
		end := start + 50

		if end > len(songsToAdd) {
			end = len(songsToAdd)
		}

		_, errAddSongs := spotifyClient.AddTracksToPlaylist(ctx, playlist.ID, songsToAdd[start:end]...)

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
