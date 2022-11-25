package song

import (
	"context"
	"log"

	"github.com/RubenPari/clear-songs/src/models"
	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/RubenPari/clear-songs/src/modules/utils"
	"github.com/gofiber/fiber/v2"
	spotifyAPI "github.com/zmb3/spotify/v2"
)

// Summary get a list of artistLibrarySummary objects
// TODO: fix
func Summary(c *fiber.Ctx) error {
	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	// call to spotify api n times based on the offset
	var limit = 50
	var offset = 0 // NOTE: offset is excluded
	// NOTE: gli elementi selezionati da n a m range
	// non sono salvati in tale ordine

	var errGetSavedTracks error

	songs := make([]spotifyAPI.SavedTrack, 0)

	for {
		tracksPage, err := spotifyClient.CurrentUsersTracks(ctx, spotifyAPI.Limit(limit), spotifyAPI.Offset(offset))

		if err != nil {
			errGetSavedTracks = err
			break
		}

		songs = append(songs, tracksPage.Tracks...)

		if len(tracksPage.Tracks) < limit {
			break
		}

		offset += limit
	}

	if errGetSavedTracks != nil {
		log.Default().Printf("couldn't get songs: %v", errGetSavedTracks)
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't get songs",
		})
	}

	artistsSummary := make([]models.ArtistLibrarySummary, 0)
	var artistsFounded [100]string

	for _, song := range songs {
		founded := utils.ArrayArtistFoundedContains(&artistsFounded, string(song.Artists[0].ID))

		if !founded {
			artistsSummary = append(artistsSummary, models.ArtistLibrarySummary{
				Id:       string(song.Artists[0].ID),
				Name:     song.Artists[0].Name,
				NumSongs: 1,
			})

			utils.
				AppendArray(&artistsFounded, string(song.Artists[0].ID))
		} else {
			founded, index := utils.ArrayArtistSummaryContains(&artistsSummary, string(song.Artists[0].ID))

			if founded {
				artistsSummary[index].NumSongs++
			}
		}
	}

	_ = c.SendStatus(200)
	return c.JSON(artistsSummary)
}

// RemoveByArtist remove all songs by artist
func RemoveByArtist(c *fiber.Ctx) error {
	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	var artistId = c.Params("id_artist")

	var limit = 50
	var offset = 0

	var errGetSavedTracks error

	var songsToRemove []spotifyAPI.ID

	for {
		tracksPage, err := spotifyClient.CurrentUsersTracks(ctx, spotifyAPI.Limit(limit), spotifyAPI.Offset(offset))

		if err != nil {
			errGetSavedTracks = err
			break
		}

		// add to songsToRemove array all songs getted
		// by specific range and offset for each api call
		for _, song := range tracksPage.Tracks {
			if string(song.Artists[0].ID) == artistId {
				songsToRemove = append(songsToRemove, song.ID)
			}
		}

		if len(tracksPage.Tracks) < limit {
			break
		}

		offset += limit
	}

	if errGetSavedTracks != nil {
		log.Default().Printf("couldn't get songs: %v", errGetSavedTracks)
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't get songs",
		})
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

	_ = c.SendStatus(200)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "songs removed",
	})
}
