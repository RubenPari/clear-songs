package routes

import (
	authCONTR "github.com/RubenPari/clear-songs/src/controllers/auth"
	playlistCONTR "github.com/RubenPari/clear-songs/src/controllers/playlist"
	songCONTR "github.com/RubenPari/clear-songs/src/controllers/song"
	utilsCONTR "github.com/RubenPari/clear-songs/src/controllers/utils"
	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	songs := app.Group("/songs")

	songs.Get("/summary", songCONTR.Summary)
	songs.Delete("/remove-by-artist/:id_artist", songCONTR.RemoveByArtist)

	songsMultiple := songs.Group("/multiple")

	songsMultiple.Delete("/remove-by-artists", songCONTR.MultipleRemoveByArtist)

	playlist := app.Group("/playlist")

	playlist.Post("/rap/create", playlistCONTR.CreateRapPlaylist)
	playlist.Post("/edm/create", playlistCONTR.CreateEdmPlaylist)

	utils := app.Group("/utils")

	utils.Get("/artist/get-id/:name", utilsCONTR.GetIdByName)
	utils.Get("/artist/get-name/:id", utilsCONTR.GetNameById)

	auth := app.Group("/auth")

	auth.Get("/login", authCONTR.Login)
	auth.Get("/callback", authCONTR.Callback)
}
