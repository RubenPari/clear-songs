package controllers

import (
	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/services"
	"github.com/gin-gonic/gin"
	spotifyAPI "github.com/zmb3/spotify"
)

func GetAll(c *gin.Context) {
	// get albums from user
	var albums []spotifyAPI.SavedAlbum

	value, found := cacheManager.Get("albums")

	if found {
		albums = value.([]spotifyAPI.SavedAlbum)
	} else {
		albums = services.GetAllUserAlbums()

		// save user albums in cacheManager
		cacheManager.Set("albums", albums)
	}

	c.JSON(200, albums)
}

func GetAlbumByArtist(c *gin.Context) {
	idArtist := spotifyAPI.ID(c.Param("id_artist"))

	// get albums from user
	var albums []spotifyAPI.SavedAlbum

	value, found := cacheManager.Get("albums")

	if found {
		albums = value.([]spotifyAPI.SavedAlbum)
	} else {
		albums = services.GetAllUserAlbums()

		// save user albums in cacheManager
		cacheManager.Set("albums", albums)
	}

	albumsArtist := services.GetAllUserAlbumsByArtist(idArtist, albums)

	c.JSON(200, albumsArtist)
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
