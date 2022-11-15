package models

import spotifyAPI "github.com/zmb3/spotify/v2"

type Artist struct {
	Id   spotifyAPI.ID  `json:"id"`
	Name string         `json:"name"`
	Uri  spotifyAPI.URI `json:"uri"`
}
