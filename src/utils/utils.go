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

func GenerateRandomWord(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rand.NewSource(time.Now().UnixNano())

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
	var genres []string

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

// GetAllUserTracks returns all
// tracks of user call the endpoint:
// https://api.spotify.com/v1/me/tracks
// with limit 50 and offset 0
// and repeat the call with offset 50, 100, 150, etc.
// until the response is empty
func GetAllUserTracks() ([]spotifyAPI.SavedTrack, error) {
	var allTracks []spotifyAPI.SavedTrack
	var offset = 0
	var limit = 50

	log.Default().Println("Getting all user tracks")

	for {
		tracks, err := SpotifyClient.CurrentUsersTracksOpt(&spotifyAPI.Options{
			Limit:  &limit,
			Offset: &offset,
		})

		log.Default().Println("Getting tracks from offset: ", offset)

		if err != nil {
			log.Default().Println("Error getting user tracks")
			return nil, err
		}

		if len(tracks.Tracks) == 0 {
			break
		}

		allTracks = append(allTracks, tracks.Tracks...)

		offset += 50
	}

	log.Println("Total tracks: ", len(allTracks))

	return allTracks, nil
}

// GetAllUserTracks returns
// all tracks of user
func GetAllUserTracksByArtist(id spotifyAPI.ID) ([]spotifyAPI.ID, error) {
	var filtredTracks []spotifyAPI.ID
	var offset = 0
	var limit = 50

	log.Default().Println("Getting all user tracks")

	for {
		tracks, err := SpotifyClient.CurrentUsersTracksOpt(&spotifyAPI.Options{
			Limit:  &limit,
			Offset: &offset,
		})

		log.Default().Println("Getting tracks from offset: ", offset)

		if err != nil {
			log.Default().Println("Error getting user tracks")
			return nil, err
		}

		if len(tracks.Tracks) == 0 {
			break
		}

		// filter by artist id
		for _, track := range tracks.Tracks {
			if track.Artists[0].ID == id {
				filtredTracks = append(filtredTracks, track.ID)
			}
		}

		offset += 50
	}

	log.Println("Total tracks: ", len(filtredTracks))

	return filtredTracks, nil
}

// DeleteTracksUser deletes
// all tracks of user
func DeleteTracksUser(tracks []spotifyAPI.ID) error {
	var offset = 0
	var limit = 50

	log.Default().Println("Deleting all user tracks")

	for {
		err := SpotifyClient.RemoveTracksFromLibrary(tracks[offset : offset+limit]...)

		log.Default().Println("Deleting tracks from offset: ", offset)

		if err != nil {
			log.Default().Println("Error deleting user tracks")
			return err
		}

		if offset > len(tracks) {
			break
		}

		offset += 50
	}

	log.Default().Println("Deleted all track")

	return nil
}

// GetAllUserTracksByGenre return all
// user track library by genre
func GetAllUserTracksByGenre(genre string) ([]spotifyAPI.ID, error) {
	// get all possible genres name
	genres := GetPossibleGenres(genre)

	var tracksFilter []spotifyAPI.ID

	var offset = 0
	var limit = 50

	log.Default().Println("Getting all user tracks by genre")

	for {
		tracks, err := SpotifyClient.CurrentUsersTracksOpt(&spotifyAPI.Options{
			Limit:  &limit,
			Offset: &offset,
		})

		log.Default().Println("Getting tracks from offset: ", offset)

		if err != nil {
			log.Default().Println("Error getting user tracks")
			return nil, err
		}

		if len(tracks.Tracks) == 0 {
			break
		}

		// filter by genre name
		for _, track := range tracks.Tracks {
			// get artist info object
			artist, _ := SpotifyClient.GetArtist(track.Artists[0].ID)

			// check if artist has the specific genre
			if ContainsGenre(artist.Genres, genres) {
				tracksFilter = append(tracksFilter, track.ID)
			}
		}

		offset += 50
	}

	return tracksFilter, nil
}

// FilterByMin returns an array of tracks
// of artist that have at least the
// minimum number of tracks
func FilterByMin(tracks map[string]int, min int) map[string]int {
	var newTracks = make(map[string]int)

	for artist, count := range tracks {
		if count >= min {
			newTracks[artist] = count
		}
	}

	return newTracks
}

// FilterByMax returns an array of tracks
// of artist that have at most the
// maximum number of tracks
func FilterByMax(tracks map[string]int, max int) map[string]int {
	var newTracks = make(map[string]int)

	for artist, count := range tracks {
		if count <= max {
			newTracks[artist] = count
		}
	}

	return newTracks
}
