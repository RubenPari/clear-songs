package controllers

import (
	"github.com/RubenPari/clear-songs/src/services"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

func GetAll(c *gin.Context) {
	albums := services.GetAllUserAlbums()

	c.JSON(200, albums)
}

func GetAlbumByArtist(c *gin.Context) {
	idArtist := spotifyAPI.ID(c.Param("id_artist"))

	albums := services.GetAllUserAlbumsByArtist(idArtist)

	c.JSON(200, albums)
}

func ConvertAlbumToSongs(c *gin.Context) {
	idAlbum := spotifyAPI.ID(c.Query("id_album"))

	errConvert := services.ConvertAlbumToSongs(idAlbum)

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
