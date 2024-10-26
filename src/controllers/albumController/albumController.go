package albumController

import (
	"github.com/RubenPari/clear-songs/src/services/userService"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

// ConvertAlbumToSongs converts an album to songs and saves them to the user library.
//
// It takes the album ID as a query parameter and returns a JSON response with a 200 status
// code and a success message if the conversion is successful. Otherwise, it returns a JSON
// response with a 400 status code and an error message.
func ConvertAlbumToSongs(c *gin.Context) {
	idAlbum := spotifyAPI.ID(c.Query("id_album"))

	errConvert := userService.ConvertAlbumToSongs(c, idAlbum)

	if errConvert != nil {
		c.JSON(400, gin.H{
			"message": "Error converting album to songs",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Album converted to songs",
	})
}
