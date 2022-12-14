package song

import (
	"log"
	"strings"

	"github.com/RubenPari/clear-songs/src/modules/client"

	"github.com/RubenPari/clear-songs/src/models"
	"github.com/RubenPari/clear-songs/src/modules/utils"
	"github.com/gofiber/fiber/v2"
	spotifyAPI "github.com/zmb3/spotify/v2"
)

// Summary get a list of artistLibrarySummary objects
// TODO: fix
func Summary(c *fiber.Ctx) error {
	// get all songs in my library
	songs, errAllTracks := utils.GetTracksUser()

	if errAllTracks != nil {
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't get all tracks",
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
	var artistId = c.Params("id_artist")

	// get all songs in my library
	allTracks, errAllTracks := utils.GetTracksUser()

	if errAllTracks != nil {
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't get all tracks",
		})
	}

	// filter songs by artist
	songsToRemove := make([]spotifyAPI.ID, 0)

	for _, track := range allTracks {
		if string(track.Artists[0].ID) == artistId {
			songsToRemove = append(songsToRemove, track.ID)
		}
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

	_ = c.SendStatus(200)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "songs removed",
	})
}

// MultipleRemoveByArtist remove all songs
// by multiple artists passed in body
func MultipleRemoveByArtist(c *fiber.Ctx) error {
	// parse body as raw/text containing a list of artists id
	artistsIdString := string(c.Body())
	artistsId := strings.Split(artistsIdString, ",")

	// call a client function to remove songs for each artistId
	for _, artistId := range artistsId {
		errRemoveByArtist := client.RemoveSongsByArtist(artistId)

		if errRemoveByArtist != nil {
			log.Default().Printf("couldn't remove songs: %v", errRemoveByArtist)
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
		"message": "songs removed for all artists",
	})
}
