package auth

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	spotifyAPI "github.com/zmb3/spotify/v2"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	spotifyOAuth "golang.org/x/oauth2/spotify"
)

var (
	SpotifyClient *spotifyAPI.Client = nil
)

// TODO: move func in utils package
// LoadEnv loads environment
// variables from .env file
// in the root directory
// upDir: number of directories
// to move up from the current path
func LoadEnv(upDir int) bool {
	// get current path
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	var numUpDir = ""

	for i := 0; i < upDir; i++ {
		numUpDir += "../"
	}

	// move up of n directories
	rootPath := filepath.Join(basePath, numUpDir)

	// load env variables
	err := godotenv.Load(rootPath + "/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return true
}

func GenerateStateString() string {
	// TODO: generate random string
	return "state"
}

func GetOAuthConfig() (*oauth2.Config, string) {
	loaded := LoadEnv(3)

	if !loaded {
		log.Default().Println("couldn't load env variables")
	}

	scopes := os.Getenv("SCOPES")
	scopesArray := strings.Split(scopes, ",")

	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       scopesArray,
		Endpoint:     spotifyOAuth.Endpoint,
	}, GenerateStateString()
}
