package utils

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/RubenPari/clear-songs/src/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	SpotifyClient *spotifyAPI.Client
	TokenHeader   string
)

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

func RandomString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, n)

	for i := range result {
		result[i] = letters[rng.Intn(len(letters))]
	}

	return string(result)
}

func GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URI"),
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

// Contains checks if an array of string
// contains an element string
func Contains(array []string, element string) bool {
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
		if Contains(genresToSearch, genre) {
			return true
		}
	}

	return false
}

// FilterSummaryByRange returns an array of
// artist summary that have at least the
// minimum number of tracks and at most the
// maximum number of tracks
// NOTE: if min or max are 0, they are ignored
func FilterSummaryByRange(tracks []models.ArtistSummary, min int, max int) []models.ArtistSummary {
	log.Default().Println("Filtering artist summary array by range")

	var newTracks []models.ArtistSummary

	for _, track := range tracks {
		if min == 0 && max == 0 {
			newTracks = append(newTracks, track)
		} else if min == 0 && track.Count <= max {
			newTracks = append(newTracks, track)
		} else if max == 0 && track.Count >= min {
			newTracks = append(newTracks, track)
		} else if track.Count >= min && track.Count <= max {
			newTracks = append(newTracks, track)
		}
	}

	return newTracks
}

func FilterGroupSummaryByRange(tracks []models.ArtistGroupSummary, min int, max int) []models.ArtistGroupSummary {
	log.Default().Println("Filtering artist summary array by range")

	var filteredGroupSummaries []models.ArtistGroupSummary

	for _, groupSummary := range tracks {
		var filteredArtists []models.ArtistSummary
		for _, artist := range groupSummary.Artists {
			if (min == 0 || artist.Count >= min) && (max == 0 || artist.Count <= max) {
				filteredArtists = append(filteredArtists, artist)
			}
		}

		// Aggiungi solo i gruppi che hanno artisti filtrati
		if len(filteredArtists) > 0 {
			filteredGroupSummary := models.ArtistGroupSummary{
				Genre:   groupSummary.Genre,
				Artists: filteredArtists,
			}
			filteredGroupSummaries = append(filteredGroupSummaries, filteredGroupSummary)
		}
	}

	return filteredGroupSummaries
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
