package database

import (
	"errors"
	"log"

	"github.com/RubenPari/clear-songs/src/domain/interfaces"
	"github.com/RubenPari/clear-songs/src/models"
	spotifyAPI "github.com/zmb3/spotify"
	"gorm.io/gorm"
)

// PostgresRepository implements DatabaseRepository interface
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new Postgres repository
// If db is nil, returns a no-op repository
func NewPostgresRepository(db *gorm.DB) interfaces.DatabaseRepository {
	if db == nil {
		return &NoOpDatabaseRepository{}
	}
	return &PostgresRepository{db: db}
}

// SaveTracksBackup saves tracks to database as backup
func (r *PostgresRepository) SaveTracksBackup(tracks []spotifyAPI.PlaylistTrack) error {
	log.Println("Saving tracks backup started")

	for _, trackPlaylist := range tracks {
		track := models.TrackDB{
			Id:     trackPlaylist.Track.ID.String(),
			Name:   trackPlaylist.Track.Name,
			Artist: trackPlaylist.Track.Artists[0].Name,
			Album:  trackPlaylist.Track.Album.Name,
			URI:    string(trackPlaylist.Track.URI),
			URL:    trackPlaylist.Track.ExternalURLs["spotify"],
		}

		var existingTrack models.TrackDB
		result := r.db.First(&existingTrack, "id = ?", track.Id)

		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Printf("Error querying existing track: %v\n", result.Error)
				return result.Error
			}

			// Track doesn't exist, insert it
			if err := r.db.Create(&track).Error; err != nil {
				log.Printf("Error inserting track: %v - %v\n", track, err)
				return err
			}
		}
		// If track exists, skip it
	}

	return nil
}

// NoOpDatabaseRepository is a no-op implementation when database is not available
type NoOpDatabaseRepository struct{}

func (n *NoOpDatabaseRepository) SaveTracksBackup(tracks []spotifyAPI.PlaylistTrack) error {
	log.Println("WARNING: Database not available, skipping track backup")
	return nil // No-op
}

// Ensure implementations
var _ interfaces.DatabaseRepository = (*PostgresRepository)(nil)
var _ interfaces.DatabaseRepository = (*NoOpDatabaseRepository)(nil)
