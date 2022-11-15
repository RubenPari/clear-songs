package song

import (
	"context"
	"log"

	trackDB "github.com/RubenPari/clear-songs/src/database/track"

	"github.com/RubenPari/clear-songs/src/models"
	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/RubenPari/clear-songs/src/modules/utils"
	"github.com/gofiber/fiber/v2"
)

// GetAllSongs get all songs of the user
// logged and save it in the database
func GetAllSongs(c *fiber.Ctx) error {
	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	songs, err := spotifyClient.CurrentUsersTracks(ctx, nil)

	if err != nil {
		log.Default().Printf("couldn't get songs: %v", err)
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't get songs",
		})
	}

	tracks := make([]models.Track, 0)

	for _, song := range songs.Tracks {
		track := models.Track{
			Id:        song.ID,
			Name:      song.Name,
			Uri:       song.URI,
			Album:     song.Album.Name,
			Artist:    song.Artists[0].Name,
			Featuring: utils.GetFeaturing(song),
		}

		tracks = append(tracks, track)
	}

	saved := trackDB.Adds(tracks)

	if !saved {
		_ = c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't save songs",
		})
	}

	_ = c.SendStatus(fiber.StatusCreated)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "songs saved",
	})
}
