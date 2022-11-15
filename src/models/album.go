package models

import spotifyAPI "github.com/zmb3/spotify/v2"

type Album struct {
	Id     spotifyAPI.ID  `json:"id"`
	Name   string         `json:"name"`
	Uri    spotifyAPI.URI `json:"uri"`
	Artist string         `json:"artist"`
}
