package client

import (
	"fmt"
	"net/http"

	"github.com/RubenPari/clear-songs/src/lib/utils"
	spotifyAPI "github.com/zmb3/spotify"
)

func RemoveAlbumsForUser(albumsId []spotifyAPI.ID) (bool, error) {
	for _, albumId := range albumsId {
		_, err := RemoveAlbumForUser(albumId)

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// RemoveAlbumForUser remove an followed album for user
// fetching spotify API without client
func RemoveAlbumForUser(albumId spotifyAPI.ID) (bool, error) {
	// get AccessToken
	accessTokenObj, _ := utils.SpotifyClient.Token()
	accessToken := accessTokenObj.AccessToken

	// create url
	url := fmt.Sprintf("https://api.spotify.com/v1/me/albums?ids=%s", albumId)

	// create request
	req, errReq := http.NewRequest("DELETE", url, nil)

	if errReq != nil {
		return false, errReq
	}

	// set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// create client
	client := &http.Client{}

	// do request
	resp, errResp := client.Do(req)

	if errResp != nil {
		return false, errResp
	}

	// check status code
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("Error removing album for user")
	}

	return true, nil
}
