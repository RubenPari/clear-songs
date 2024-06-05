package utils

import (
	"encoding/json"
	"errors"
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io"
	"log"
	"os"
)

var SpotifyClient *spotifyAPI.Client

// ClientID ClientSecret RedirectURI Port
const ClientID = "06d2f7ccaabd48829ad97f299c13c1be"
const ClientSecret = "ecc19973c7d7459fa2fd6a4206ae538a"
const RedirectURI = "http://localhost:3000/auth/callback"
const Port = "3000"
const fileNameTracksBackup = "tracks-backup.json"
const FilePermission = 0755

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

// SaveTracksFileBackupIDs saves a list of track IDs
// to a json file in the root directory for recovery
// track in case of accidental deletion
func SaveTracksFileBackupIDs(tracksIds []spotifyAPI.ID) error {
	// open file for reading and writing if it doesn't exist
	file, errOpenFile := os.OpenFile(fileNameTracksBackup, os.O_RDWR|os.O_CREATE, FilePermission)

	if errOpenFile != nil {
		return errOpenFile
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Read existing IDs from file
	var existingTrackIDs []spotifyAPI.ID
	decoder := json.NewDecoder(file)
	errReadIDs := decoder.Decode(&existingTrackIDs)

	if errReadIDs != nil && errOpenFile != io.EOF {
		return errReadIDs
	}

	// Create a map of existing IDs for faster lookup
	idMap := make(map[spotifyAPI.ID]bool)
	for _, id := range existingTrackIDs {
		idMap[id] = true
	}

	// Add new IDs to the existing IDs
	for _, trackId := range tracksIds {
		_, exists := idMap[trackId]

		if !exists {
			existingTrackIDs = append(existingTrackIDs, trackId)
			idMap[trackId] = true
		}
	}

	// Truncate the file to 0
	errTruncateFile := file.Truncate(0)

	if errTruncateFile != nil {
		return errTruncateFile
	}

	// Write the new IDs to the file
	encoder := json.NewEncoder(file)
	errWriteTrackIDs := encoder.Encode(existingTrackIDs)

	if errWriteTrackIDs != nil {
		return errWriteTrackIDs
	}

	return nil
}
