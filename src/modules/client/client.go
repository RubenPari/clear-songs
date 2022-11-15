package client

import (
	"encoding/json"
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
