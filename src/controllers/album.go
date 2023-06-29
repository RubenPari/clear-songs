package controllers

import (
	userService "github.com/RubenPari/clear-songs/src/services/user"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

func GetAll(c *gin.Context) {
	albums := userService.GetAllUserAlbums()

	c.JSON(200, albums)
}

func GetAlbumByArtist(c *gin.Context) {
	idArtist := spotifyAPI.ID(c.Param("id_artist"))

	albums := userService.GetAllUserAlbumsByArtist(idArtist)

	c.JSON(200, albums)
}

func DeleteAlbumByArtist(c *gin.Context) {
	idArtist := spotifyAPI.ID(c.Param("id_artist"))

	err := userService.DeleteAlbumsByArtist(idArtist)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Error deleting albums",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Albums deleted",
	})
}

func ConvertAlbumToSongs(c *gin.Context) {
	idAlbum := spotifyAPI.ID(c.Query("id_album"))

	err := userService.ConvertAlbumToSongs(idAlbum)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Error converting album to songs",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Album converted to songs",
	})
}
