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
