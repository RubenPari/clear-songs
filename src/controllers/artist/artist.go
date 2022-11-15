package artist

import (
	"context"
	"log"
	"net/http"

	trackDB "github.com/RubenPari/clear-songs/src/database/track"

	artistDB "github.com/RubenPari/clear-songs/src/database/artist"
	"github.com/RubenPari/clear-songs/src/models"
	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/gofiber/fiber/v2"
	spotifyAPI "github.com/zmb3/spotify/v2"
)

func Add(c *fiber.Ctx) error {
	id := c.Params("id")

	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	artistApi, err := spotifyClient.GetArtist(ctx, spotifyAPI.ID(id))

	if err != nil || artistApi == nil {
		log.Fatalf("couldn't get artist information: %v", err)
	}

	artistObj := models.Artist{
		Id:   artistApi.ID,
		Name: artistApi.Name,
		Uri:  artistApi.URI,
	}

	success := artistDB.Add(artistObj)

	if success {
		_ = c.SendStatus(http.StatusCreated)
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "artist added",
		})
	} else {
		_ = c.SendStatus(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "error adding artist",
		})
	}
}

func CheckIfSaved(c *fiber.Ctx) error {
	id := c.Params("id")

	artistObj := artistDB.Get(id)

	if artistObj.Id == ("") {
		_ = c.SendStatus(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "artist not found",
		})
	}

	_ = c.SendStatus(http.StatusOK)
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "artist found",
	})
}

// GetAllSongs get all songs of
// a specific artist
// and add them to the database
func GetAllSongs(c *fiber.Ctx) error {
	id := c.Params("id")

	spotifyClient := authMO.SpotifyClient
	ctx := context.Background()

	artistObj := artistDB.Get(id)

	if artistObj.Id == ("") {
		_ = c.SendStatus(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "artist not found",
		})
	}

	albumsObj := make([]models.Album, 0)
	tracksObj := make([]models.Track, 0)

	// get all albums
	albums, err := spotifyClient.GetArtistAlbums(
		ctx,
		artistObj.Id,
		[]spotifyAPI.AlbumType{
			spotifyAPI.AlbumTypeAlbum,
			spotifyAPI.AlbumTypeSingle,
			spotifyAPI.AlbumTypeCompilation,
		})

	if err != nil {
		log.Fatalf("couldn't get artist albums: %v", err)
	}

	for _, album := range albums.Albums {
		albumObj := models.Album{
			Id:     album.ID,
			Name:   album.Name,
			Uri:    album.URI,
			Artist: album.Artists[0].Name,
		}

		albumsObj = append(albumsObj, albumObj)
	}

	// get all songs
	for _, album := range albumsObj {
		tracks, err := spotifyClient.GetAlbumTracks(ctx, album.Id)

		if err != nil {
			log.Fatalf("couldn't get album tracks: %v", err)
		}

		for _, song := range tracks.Tracks {
			featuringArtists := ""
			for i := 1; i < len(song.Artists); i++ {
				featuringArtists += song.Artists[i].Name + ", "
			}

			// check if first artist is the artist
			// we are looking for
			if song.Artists[0].Name == artistObj.Name {
				log.Default().Println("first artist is the artist we are looking for")

				trackObj := models.Track{
					Id:        song.ID,
					Name:      song.Name,
					Uri:       song.URI,
					Album:     album.Name,
					Artist:    song.Artists[0].Name,
					Featuring: featuringArtists,
				}

				tracksObj = append(tracksObj, trackObj)
			}
		}
	}

	// save songs to db
	success := trackDB.Adds(tracksObj)

	if success {
		_ = c.SendStatus(http.StatusCreated)
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "songs added",
		})
	} else {
		_ = c.SendStatus(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "error adding songs",
		})
	}
}

// GetFeaturedSongs filters all songs of
// a specific artist
// where Eminem is featured
// and add them to the database
func GetFeaturedSongs(c *fiber.Ctx) error {
	id := c.Params("id")

	tracksObj := trackDB.GetAllByArtist(id)

	tracksFiltersObj := trackDB.FilterByFeaturing(tracksObj)

	success := trackDB.AddsFeatured(tracksFiltersObj)

	if success {
		_ = c.SendStatus(http.StatusCreated)
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "songs added",
		})
	} else {
		_ = c.SendStatus(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "error adding songs",
		})
	}
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	artistObj := artistDB.Get(id)

	if artistObj.Id == ("") {
		_ = c.SendStatus(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "artist not found",
		})
	}

	artistDeleted := artistDB.Delete(id)

	// delete all songs of the artist in tracks_feat
	featuringDeleted := trackDB.DeleteAllByFeaturing(id)

	// delete all songs of the artist in tracks
	tracksDeleted := trackDB.DeleteAllByArtist(id)

	if artistDeleted && featuringDeleted && tracksDeleted {
		_ = c.SendStatus(http.StatusOK)
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "artist deleted",
		})
	} else {
		_ = c.SendStatus(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "error deleting artist",
		})
	}
}
