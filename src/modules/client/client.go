package client

import (
	"encoding/json"
	"errors"
	authMO "github.com/RubenPari/clear-songs/src/modules/auth"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

// GetNameArtistById
// call to endpoint to get name of artist by id
func GetNameArtistById(id string) (string, error) {
	_ = godotenv.Load()
	port := os.Getenv("PORT")

	// get name of artist by id with endpoint
	responseName, errGetName := http.Get("http://localhost:" + port + "/utils/artist/get-name/" + id)
	if errGetName != nil {
		log.Default().Println("Error getting name of artist")
		log.Default().Println(errGetName)
		return "", errGetName
	}

	// extract name of artist from response of type json
	type Response struct {
		Status string `json:"status"`
		Name   string `json:"name"`
	}

	var response Response
	_ = json.NewDecoder(responseName.Body).Decode(&response)

	return response.Name, nil
}

// RemoveSongsByArtist
// call to endpoint to remove songs by artist
func RemoveSongsByArtist(id string) error {
	_ = godotenv.Load()
	port := os.Getenv("PORT")

	// remove songs by artist with endpoint DELETE
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://localhost:"+port+"/songs/remove-by-artist/"+id, nil)
	req.Header.Add("Bearer", authMO.Token.AccessToken)
	resp, err := client.Do(req)

	if err != nil {
		log.Default().Println("Error removing songs by artist")
		log.Default().Println(err)
		return err
	}

	if resp.StatusCode != 200 {
		log.Default().Println("Error removing songs by artist")
		log.Default().Println(resp)
		return errors.New("error removing songs by artist")
	}

	_ = resp.Body.Close()

	return nil
}
