package postgres

import (
	"context"
	"errors"

	"github.com/RubenPari/clear-songs/internal/domain/auth"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() auth.UserRepository {
	return &userRepository{
		db: Db, // using global Db from postgres package
	}
}

func mapToUser(dbUser *models.UserDB) *auth.User {
	if dbUser == nil {
		return nil
	}
	return &auth.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		IsVerified:   dbUser.IsVerified,
		SpotifyID:    dbUser.SpotifyID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

func (r *userRepository) Create(ctx context.Context, user *auth.User) error {
	dbUser := &models.UserDB{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
		SpotifyID:    user.SpotifyID,
	}

	result := r.db.WithContext(ctx).Create(dbUser)
	if result.Error != nil {
		return result.Error
	}

	// Map generated ID and times back
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*auth.User, error) {
	var dbUser models.UserDB
	result := r.db.WithContext(ctx).First(&dbUser, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found to simplify checks
		}
		return nil, result.Error
	}
	return mapToUser(&dbUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*auth.User, error) {
	var dbUser models.UserDB
	result := r.db.WithContext(ctx).First(&dbUser, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return mapToUser(&dbUser), nil
}

func (r *userRepository) GetBySpotifyID(ctx context.Context, spotifyID string) (*auth.User, error) {
	var dbUser models.UserDB
	result := r.db.WithContext(ctx).First(&dbUser, "spotify_id = ?", spotifyID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return mapToUser(&dbUser), nil
}

func (r *userRepository) Update(ctx context.Context, user *auth.User) error {
	dbUser := &models.UserDB{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
		SpotifyID:    user.SpotifyID,
	}

	result := r.db.WithContext(ctx).Save(dbUser)
	return result.Error
}
