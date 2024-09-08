package utils

import (
	"errors"
	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
	"gorm.io/gorm"
	"log"
)

var SpotifyClient spotifyAPI.Client

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
	for _, trackPlaylist := range tracksPlaylist {
		track := models.TrackDB{
			Id:     trackPlaylist.Track.ID.String(),
			Name:   trackPlaylist.Track.Name,
			Artist: trackPlaylist.Track.Artists[0].Name,
			Album:  trackPlaylist.Track.Album.Name,
			URI:    string(trackPlaylist.Track.URI),
			URL:    trackPlaylist.Track.ExternalURLs["spotify"],
		}

		var existingTrack models.TrackDB
		alreadyExistTrack := database.Db.First(&existingTrack, "id = ?", track.Id)

		if alreadyExistTrack != nil {
			if !errors.Is(alreadyExistTrack.Error, gorm.ErrRecordNotFound) {
				log.Printf("Error querying track: %v\n", alreadyExistTrack)
				return alreadyExistTrack.Error
			}

			insertTrack := database.Db.Create(&track)

			if insertTrack.Error != nil {
				log.Printf("Error inserting track: %v\n", insertTrack.Error)
				return insertTrack.Error
			}
		}
	}

	return nil
}
