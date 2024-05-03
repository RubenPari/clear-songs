package utils

import (
	"github.com/RubenPari/clear-songs/src/models"
	"github.com/joho/godotenv"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"log"
	"os"
	"path/filepath"
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
