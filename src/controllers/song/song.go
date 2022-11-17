package song

import (
	"context"
	"github.com/RubenPari/clear-songs/src/models"
	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/RubenPari/clear-songs/src/modules/utils"
	"github.com/gofiber/fiber/v2"
	spotifyAPI "github.com/zmb3/spotify/v2"
	"log"
)

// Summary get a list of artistLibrarySummary objects
func Summary(c *fiber.Ctx) error {
	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	songs, err := spotifyClient.CurrentUsersTracks(ctx, spotifyAPI.Limit(10))

	if err != nil {
		log.Default().Printf("couldn't get songs: %v", err)
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't get songs",
		})
	}
	artistsSummary := make([]models.ArtistLibrarySummary, 0)
	var artistsFounded [50]string

	for _, song := range songs.Tracks {
		founded, index := utils.ArrayContains(artistsFounded, string(song.Artists[0].ID))

		if !founded {
			artistsSummary = append(artistsSummary, models.ArtistLibrarySummary{
				Id:   string(song.Artists[0].ID),
				Name: song.Artists[0].Name,
				Num:  1,
			})

			utils.AppendArray(&artistsFounded, string(song.Artists[0].ID))
		} else {
			artistsSummary[index].Num++
		}
	}

	_ = c.SendStatus(200)
	return c.JSON(artistsSummary)
}
