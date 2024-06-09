package utils

import (
	"errors"
	"log"

	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var SpotifyClient *spotifyAPI.Client

// ClientID ClientSecret RedirectURI Port
const ClientID = "06d2f7ccaabd48829ad97f299c13c1be"
const ClientSecret = "ecc19973c7d7459fa2fd6a4206ae538a"
const RedirectURI = "http://localhost:3000/auth/callback"
const Port = "3000"

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

func SaveTracksBackup(tracksPlaylist []spotifyAPI.PlaylistTrack) error {
	db := database.GetDB()

	for _, trackPlaylist := range tracksPlaylist {
		track := models.Track{
			Id:     trackPlaylist.Track.ID.String(),
			Name:   trackPlaylist.Track.Name,
			Artist: trackPlaylist.Track.Artists[0].Name,
			Album:  trackPlaylist.Track.Album.Name,
			URI:    string(trackPlaylist.Track.URI),
			URL:    trackPlaylist.Track.ExternalURLs["spotify"],
		}

		var existingTrack models.Track
		errAlreadyExistTRack := db.QueryRow("SELECT * FROM tracks WHERE id = ?", track.Id).Scan(&existingTrack)

		if errAlreadyExistTRack != nil && !errors.Is(errAlreadyExistTRack, gorm.ErrRecordNotFound) {
			log.Printf("Error querying track: %v\n", errAlreadyExistTRack)
			return errAlreadyExistTRack
		}

		if errors.Is(errAlreadyExistTRack, gorm.ErrRecordNotFound) {
			_, errInsertTrack := db.Exec("INSERT INTO tracks (id, name, artist, album, uri, url) VALUES (?, ?, ?, ?, ?, ?)", track.Id, track.Name, track.Artist, track.Album, track.URI, track.URL)

			if errInsertTrack != nil {
				log.Printf("Error inserting track: %v", errInsertTrack)
				return errInsertTrack
			}
		}
	}

	return nil
}
