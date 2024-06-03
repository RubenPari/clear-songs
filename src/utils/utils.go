package utils

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var SpotifyClient *spotifyAPI.Client

// ClientID ClientSecret RedirectURI Port
var ClientID = "06d2f7ccaabd48829ad97f299c13c1be"
var ClientSecret = "ecc19973c7d7459fa2fd6a4206ae538a"
var RedirectURI = "http://localhost:3000/auth/callback"
var Port = "3000"

func GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		RedirectURL:  RedirectURI,
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
		return nil, errors.New("type not supported")
	}

	return trackIDs, nil
}

// SaveTracksFileIDs saves a list of track IDs
// to a file json in the root of the project
func SaveTracksFileIDs(ids []spotifyAPI.ID) error {
	file, errCreateFile := os.Create("tracks-backup.json")

	if errCreateFile != nil {
		return errCreateFile
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := json.NewEncoder(file)

	errEncoder := encoder.Encode(ids)

	if errEncoder != nil {
		return errEncoder
	}

	return nil
}
