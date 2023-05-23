package utils

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/RubenPari/clear-songs/src/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var SpotifyClient *spotifyAPI.Client

func LoadEnv(moveUp int) {
	// add /env or \env to the path
	// depending on the OS
	var envName string

	if os.PathSeparator == '/' {
		envName = "/.env"
	} else {
		envName = "\\.env"
	}

	// get current directory
	// and move up to the root directory
	currentDir, _ := os.Getwd()

	for i := 0; i < moveUp; i++ {
		currentDir = filepath.Dir(currentDir)
	}

	// load .env file
	err := godotenv.Load(currentDir + envName)

	log.Default().Println("Loading env file in " + currentDir + envName)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"user-read-private",
			"user-read-email",
			"user-library-read",
			"user-library-modify",
			"playlist-read-private",
			"playlist-read-collaborative",
			"playlist-modify-public",
			"playlist-modify-private",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotifyAPI.AuthURL,
			TokenURL: spotifyAPI.TokenURL,
		}}
}

// GetPossibleGenres returns an array of genres
// with all possible genres name alternatives
func GetPossibleGenres(genre string) []string {
	var genres []string

	switch genre {
	case "rock":
		genres = []string{"rock", "rock and roll", "rock & roll"}
	case "pop":
		genres = []string{"pop", "pop music"}
	case "hip-hop":
		genres = []string{"hip-hop", "hip hop", "rap", "hip hop music", "rapper", "rap music", "hip-hop music"}
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

// CheckTypeObject checks if the
// object entity spotify is of the type specified
func CheckTypeObject(typeObject string) bool {
	if typeObject == "artist" ||
		typeObject == "track" ||
		typeObject == "album" ||
		typeObject == "playlist" {
		return true
	}

	return false
}

// GetObjectName returns the name of the
// object entity spotify given the type and id
func GetObjectName(typeObject string, id string) string {
	switch typeObject {
	case "artist":
		artist, _ := SpotifyClient.GetArtist(spotifyAPI.ID(id))

		return artist.Name + " - " + typeObject

	case "album":
		album, _ := SpotifyClient.GetAlbum(spotifyAPI.ID(id))

		return album.Name + " - " + typeObject

	case "track":
		track, _ := SpotifyClient.GetTrack(spotifyAPI.ID(id))

		return track.Name + " - " + typeObject
	case "playlist":
		playlist, _ := SpotifyClient.GetPlaylist(spotifyAPI.ID(id))

		return playlist.Name + " - " + typeObject

	default:
		return ""
	}
}

// GetIDByName returns the id of the
// object entity spotify given name
func GetIDByName(name string, typeObject string) spotifyAPI.ID {
	switch typeObject {
	case "artist":
		artist, _ := SpotifyClient.Search(name, spotifyAPI.SearchTypeArtist)

		return artist.Artists.Artists[0].ID

	case "album":
		album, _ := SpotifyClient.Search(name, spotifyAPI.SearchTypeAlbum)

		return album.Albums.Albums[0].ID

	case "track":
		track, _ := SpotifyClient.Search(name, spotifyAPI.SearchTypeTrack)

		return track.Tracks.Tracks[0].ID

	case "playlist":
		playlist, _ := SpotifyClient.Search(name, spotifyAPI.SearchTypePlaylist)

		return playlist.Playlists.Playlists[0].ID

	default:
		return ""
	}
}

func SaveSummaryToFile(summary []models.ArtistSummary) error {
	file, errCreate := os.Create("summary.txt")

	if errCreate != nil {
		return errCreate
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// for each artist in the summary
	// write the name and the count
	for _, summary := range summary {
		_, errWrite := file.WriteString(summary.Name + " - " + strconv.Itoa(summary.Count) + "\n")

		if errWrite != nil {
			return errWrite
		}
	}

	return nil
}

func CorsConfig(server *gin.Engine) {
	server.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		AllowFiles: true,
		MaxAge:     86400,
	}))
}
