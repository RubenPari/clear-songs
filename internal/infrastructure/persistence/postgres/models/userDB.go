package models

import (
	"time"

	"gorm.io/gorm"
)

type UserDB struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string         `gorm:"uniqueIndex;not null"`
	PasswordHash string         `gorm:"not null"`
	IsVerified   bool           `gorm:"default:false"`
	SpotifyID    *string        `gorm:"uniqueIndex"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (UserDB) TableName() string {
	return "users"
}

type VerificationTokenDB struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `gorm:"index;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User UserDB `gorm:"foreignKey:UserID"`
}

func (VerificationTokenDB) TableName() string {
	return "verification_tokens"
}

type ResetTokenDB struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `gorm:"index;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User UserDB `gorm:"foreignKey:UserID"`
}

func (ResetTokenDB) TableName() string {
	return "reset_tokens"
}
