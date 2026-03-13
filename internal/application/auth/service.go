package auth

import (
	"context"
	"errors"
	"time"

	domainAuth "github.com/RubenPari/clear-songs/internal/domain/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) error
	ConfirmEmail(ctx context.Context, token string) error
	Login(ctx context.Context, req LoginRequest) (*domainAuth.User, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error
}

type authService struct {
	userRepo  domainAuth.UserRepository
	tokenRepo domainAuth.TokenRepository
	emailSvc  domainAuth.EmailService
}

func NewAuthService(ur domainAuth.UserRepository, tr domainAuth.TokenRepository, es domainAuth.EmailService) AuthService {
	return &authService{
		userRepo:  ur,
		tokenRepo: tr,
		emailSvc:  es,
	}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) error {
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domainAuth.User{
		Email:        req.Email,
		PasswordHash: string(hashed),
		IsVerified:   false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	tokenStr := uuid.New().String()
	vToken := &domainAuth.VerificationToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.tokenRepo.CreateVerificationToken(ctx, vToken); err != nil {
		return err
	}

	return s.emailSvc.SendVerificationEmail(ctx, user.Email, tokenStr)
}

func (s *authService) ConfirmEmail(ctx context.Context, token string) error {
	vToken, err := s.tokenRepo.GetVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if vToken == nil || vToken.ExpiresAt.Before(time.Now()) {
		return ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, vToken.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.IsVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.tokenRepo.DeleteVerificationToken(ctx, token)
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*domainAuth.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, ErrEmailNotVerified
	}

	return user, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		// Do not leak existence, just return nil
		return nil
	}

	tokenStr := uuid.New().String()
	rToken := &domainAuth.ResetToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.tokenRepo.CreateResetToken(ctx, rToken); err != nil {
		return err
	}

	return s.emailSvc.SendPasswordResetEmail(ctx, user.Email, tokenStr)
}

func (s *authService) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	rToken, err := s.tokenRepo.GetResetToken(ctx, req.Token)
	if err != nil {
		return err
	}
	if rToken == nil || rToken.ExpiresAt.Before(time.Now()) {
		return ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, rToken.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashed)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.tokenRepo.DeleteResetToken(ctx, req.Token)
}

func (s *authService) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashed)
	return s.userRepo.Update(ctx, user)
}
