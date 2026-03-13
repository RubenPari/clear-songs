package postgres

import (
	"context"
	"errors"

	"github.com/RubenPari/clear-songs/internal/domain/auth"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres/models"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository() auth.TokenRepository {
	return &tokenRepository{
		db: Db,
	}
}

func (r *tokenRepository) CreateVerificationToken(ctx context.Context, token *auth.VerificationToken) error {
	dbToken := &models.VerificationTokenDB{
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
	}

	result := r.db.WithContext(ctx).Create(dbToken)
	if result.Error != nil {
		return result.Error
	}

	token.ID = dbToken.ID
	token.CreatedAt = dbToken.CreatedAt
	return nil
}

func (r *tokenRepository) GetVerificationToken(ctx context.Context, tokenStr string) (*auth.VerificationToken, error) {
	var dbToken models.VerificationTokenDB
	result := r.db.WithContext(ctx).First(&dbToken, "token = ?", tokenStr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &auth.VerificationToken{
		ID:        dbToken.ID,
		UserID:    dbToken.UserID,
		Token:     dbToken.Token,
		ExpiresAt: dbToken.ExpiresAt,
		CreatedAt: dbToken.CreatedAt,
	}, nil
}

func (r *tokenRepository) DeleteVerificationToken(ctx context.Context, tokenStr string) error {
	result := r.db.WithContext(ctx).Where("token = ?", tokenStr).Delete(&models.VerificationTokenDB{})
	return result.Error
}

func (r *tokenRepository) CreateResetToken(ctx context.Context, token *auth.ResetToken) error {
	dbToken := &models.ResetTokenDB{
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
	}

	result := r.db.WithContext(ctx).Create(dbToken)
	if result.Error != nil {
		return result.Error
	}

	token.ID = dbToken.ID
	token.CreatedAt = dbToken.CreatedAt
	return nil
}

func (r *tokenRepository) GetResetToken(ctx context.Context, tokenStr string) (*auth.ResetToken, error) {
	var dbToken models.ResetTokenDB
	result := r.db.WithContext(ctx).First(&dbToken, "token = ?", tokenStr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &auth.ResetToken{
		ID:        dbToken.ID,
		UserID:    dbToken.UserID,
		Token:     dbToken.Token,
		ExpiresAt: dbToken.ExpiresAt,
		CreatedAt: dbToken.CreatedAt,
	}, nil
}

func (r *tokenRepository) DeleteResetToken(ctx context.Context, tokenStr string) error {
	result := r.db.WithContext(ctx).Where("token = ?", tokenStr).Delete(&models.ResetTokenDB{})
	return result.Error
}
