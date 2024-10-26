package utils

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RubenPari/clear-songs/src/services/SpotifyService"
	"github.com/gin-gonic/gin"

	"github.com/RubenPari/clear-songs/src/constants"

	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	configAuth = GetOAuth2Config()
	SpotifySvc = SpotifyService.NewSpotifyService(configAuth.ClientID, configAuth.ClientSecret, configAuth.RedirectURL)
)

// GetOAuth2Config returns a pointer to an oauth2.Config with the client id, client
// secret, redirect url, and scopes set from the environment variables CLIENT_ID,
// CLIENT_SECRET, REDIRECT_URL, and SPOTIFY_SCOPES, respectively. The endpoint is set
// to the Spotify endpoints.
func GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       constants.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotifyAPI.AuthURL,
			TokenURL: spotifyAPI.TokenURL,
		}}
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

// ConvertTracksToID converts a list of tracks
// can be of type:
// - []spotifyAPI.FullTrack,
// - []spotifyAPI.PlaylistTrack,
// - []spotifyAPI.SavedTrack,
// - []spotifyAPI.SavedAlbum
// to a list of track IDs
func ConvertTracksToID(tracks interface{}) ([]spotifyAPI.ID, error) {
	var trackIDs []spotifyAPI.ID

	switch t := tracks.(type) {
	case []spotifyAPI.FullTrack:
		for _, track := range t {
			trackIDs = append(trackIDs, track.ID)
		}
	case []spotifyAPI.PlaylistTrack:
		for _, track := range t {
			trackIDs = append(trackIDs, track.Track.ID)
		}
	case []spotifyAPI.SavedTrack:
		for _, track := range t {
			trackIDs = append(trackIDs, track.FullTrack.ID)
		}
	case []spotifyAPI.SavedAlbum:
		for _, album := range t {
			for _, track := range album.Tracks.Tracks {
				trackIDs = append(trackIDs, track.ID)
			}
		}
	default:
		return nil, errors.New(" ConvertTracksToID: Type input not supported")
	}

	return trackIDs, nil
}

// SaveTracksBackup saves a list of tracks
// in a SQLite database as a backup
//
// The tracks are saved in the `tracks` table
// with the following columns:
// - `id`: the track's ID as a string
// - `name`: the track's name
// - `artist`: the track's artist
// - `album`: the track's album
// - `uri`: the track's URI as a string
// - `url`: the track's URL on Spotify
//
// If a track already exists in the database,
// it is not inserted again
//
// If an error occurs while saving the tracks,
// it is returned as an error
func SaveTracksBackup(tracksPlaylist []spotifyAPI.PlaylistTrack) error {
	log.Default().Println("Saving tracks backup started")

	for _, trackPlaylist := range tracksPlaylist {
		track := models.TrackDB{
			Id:     trackPlaylist.Track.ID.String(),
			Name:   trackPlaylist.Track.Name,
			Artist: trackPlaylist.Track.Artists[0].Name,
			Album:  trackPlaylist.Track.Album.Name,
			URI:    string(trackPlaylist.Track.URI),
			URL:    trackPlaylist.Track.ExternalURLs["spotify"],
		}

		log.Default().Printf("Created TrackDB: Name: %s, Artist: %s\n", track.Name, track.Artist)

		var existingTrack models.TrackDB
		alreadyExistTrack := database.Db.First(&existingTrack, "id = ?", track.Id)

		if alreadyExistTrack != nil {
			if !errors.Is(alreadyExistTrack.Error, gorm.ErrRecordNotFound) {
				log.Printf("Error querying alreadyExistTrack: %v\n", alreadyExistTrack)
				return alreadyExistTrack.Error
			}

			insertTrack := database.Db.Create(&track)

			if insertTrack.Error != nil {
				log.Printf("Error inserting track: %v - %v\n", track, insertTrack.Error)
				return insertTrack.Error
			}
		}
	}

	return nil
}

// GetSpotifyService retrieves the SpotifyService from the gin context.
// It returns nil if the SpotifyService is not found in the context.
func GetSpotifyService(c *gin.Context) *SpotifyService.SpotifyService {
	service, exists := c.Get("spotifyService")
	if !exists {
		return nil
	}
	return service.(*SpotifyService.SpotifyService)
}

// LoadEnvVariables load environment variables from a file path
func LoadEnvVariables() {
	// get current working directory
	cwd, errCwd := os.Getwd()

	if errCwd != nil {
		log.Fatalf("error getting current working directory: %v", errCwd)
	}

	// move up one level folder
	cwd = filepath.Dir(cwd)

	envPath := filepath.Join(cwd, ".env")

	file, errOpenFile := os.Open(envPath)

	if errOpenFile != nil {
		log.Fatalf("error opening .env file: %v", errOpenFile)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)

	// read the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// split the line into key and value
		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			log.Fatalf("invalid line in .env file: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// remove quotes "" from the value
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			value = strings.Trim(value, `"`)
		}

		// set the environment variable
		errSetEnvVar := os.Setenv(key, value)

		if errSetEnvVar != nil {
			log.Fatalf("error setting environment variable: %v", errSetEnvVar)
		}
	}

	errReadFile := scanner.Err()

	if errReadFile != nil {
		log.Fatalf("error reading .env file: %v", errReadFile)
	}
}
