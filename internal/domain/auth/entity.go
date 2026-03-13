package auth

import (
	"context"
	"time"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	IsVerified   bool
	SpotifyID    *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type VerificationToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type ResetToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetBySpotifyID(ctx context.Context, spotifyID string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type TokenRepository interface {
	CreateVerificationToken(ctx context.Context, token *VerificationToken) error
	GetVerificationToken(ctx context.Context, token string) (*VerificationToken, error)
	DeleteVerificationToken(ctx context.Context, token string) error

	CreateResetToken(ctx context.Context, token *ResetToken) error
	GetResetToken(ctx context.Context, token string) (*ResetToken, error)
	DeleteResetToken(ctx context.Context, token string) error
}

type EmailService interface {
	SendVerificationEmail(ctx context.Context, email string, token string) error
	SendPasswordResetEmail(ctx context.Context, email string, token string) error
}
