package utils

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var SpotifyClient *spotifyAPI.Client

func LoadEnv(moveUp int) {
	// get current directory
	// and move up to the root directory
	currentDir, _ := os.Getwd()

	for i := 0; i < moveUp; i++ {
		currentDir = filepath.Dir(currentDir)
	}

	// load .env file
	err := godotenv.Load(currentDir + "/.env")

	log.Default().Println(currentDir + "/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GenerateRandomWord(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())

	word := make([]rune, n) // initialize word with the correct length

	for i := range word {
		word[i] = letters[rand.Intn(len(letters))]
	}

	return string(word)
}

func GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"user-read-private",
			"user-read-email",
			"playlist-read-private",
			"playlist-read-collaborative",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotifyAPI.AuthURL,
			TokenURL: spotifyAPI.TokenURL,
		}}
}

// ArrayContains checks if an array of string
// contains an element string
func ArrayContains(array []string, element string) bool {
	for _, a := range array {
		if a == element {
			return true
		}
	}

	return false
}

// ContainsGenre checks if an array of string genres
// contains almost one genre of the second array of genres
func ContainsGenre(genres []string, genresToSearch []string) bool {
	for _, genre := range genres {
		if ArrayContains(genresToSearch, genre) {
			return true
		}
	}

	return false
}

// GetPossibleGenres returns an array of genres
// with all possible genres name alternatives
func GetPossibleGenres(genre string) []string {
	var genres = []string{}

	switch genre {
	case "rock":
		genres = []string{"rock", "rock and roll", "rock & roll"}
	case "pop":
		genres = []string{"pop", "pop music"}
	case "hip-hop":
		genres = []string{"hip-hop", "hip hop", "rap", "hip hop music", "rappeur", "rap music", "hip-hop music"}
	case "r&b":
		genres = []string{"r&b", "rnb", "r&b music", "rnb music"}
	case "country":
		genres = []string{"country", "country music"}
	case "jazz":
		genres = []string{"jazz", "jazz music"}
	case "blues":
		genres = []string{"blues", "blues music"}
	case "metal":
		genres = []string{"metal", "metal music"}
	case "classical":
		genres = []string{"classical", "classical music"}
	case "reggae":
		genres = []string{"reggae", "reggae music"}
	case "soul":
		genres = []string{"soul", "soul music"}
	case "electronic":
		genres = []string{"electronic", "electronic music", "electro", "EDM", "electro music", "EDM music"}
	case "folk":
		genres = []string{"folk", "folk music"}
	}

	return genres
}
