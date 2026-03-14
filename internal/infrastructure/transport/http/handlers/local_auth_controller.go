package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RubenPari/clear-songs/internal/application/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LocalAuthController struct {
	BaseController
	authService auth.AuthService
}

func NewLocalAuthController(authSvc auth.AuthService) *LocalAuthController {
	return &LocalAuthController{
		authService: authSvc,
	}
}

func (ac *LocalAuthController) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.JSONValidationError(c, "Invalid request payload")
		return
	}

	ctx := context.Background()
	err := ac.authService.Register(ctx, req)
	if err != nil {
		if err == auth.ErrUserExists {
			ac.JSONValidationError(c, "User already exists")
			return
		}
		log.Printf("ERROR: Registration failed: %v", err)
		ac.JSONInternalError(c, "Registration failed")
		return
	}

	ac.JSONSuccess(c, gin.H{"message": "Registration successful. Please check your email to confirm your account."})
}

func (ac *LocalAuthController) ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		ac.JSONValidationError(c, "Token is required")
		return
	}

	ctx := context.Background()
	err := ac.authService.ConfirmEmail(ctx, token)
	if err != nil {
		ac.JSONError(c, http.StatusBadRequest, "INVALID_TOKEN", "Invalid or expired token")
		return
	}

	ac.JSONSuccess(c, gin.H{"message": "Email confirmed successfully"})
}

func (ac *LocalAuthController) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.JSONValidationError(c, "Invalid request payload")
		return
	}

	ctx := context.Background()
	user, err := ac.authService.Login(ctx, req)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			ac.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials")
			return
		}
		if err == auth.ErrEmailNotVerified {
			ac.JSONError(c, http.StatusForbidden, "FORBIDDEN", "Email not verified")
			return
		}
		ac.JSONInternalError(c, "Login failed")
		return
	}

	// Generate JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "fallback-secret-for-dev"
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		ac.JSONInternalError(c, "Failed to generate token")
		return
	}

	// Set as HTTP-only cookie
	c.SetCookie("auth_token", tokenString, 7*24*3600, "/", "", false, true)

	ac.JSONSuccess(c, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"spotify_id": user.SpotifyID,
		},
	})
}

func (ac *LocalAuthController) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.JSONValidationError(c, "Invalid request payload")
		return
	}

	ctx := context.Background()
	// Ignore errors to prevent email enumeration, but log them for debugging
	if err := ac.authService.ForgotPassword(ctx, req.Email); err != nil {
		// Log the error but don't expose it to the client
		log.Printf("ERROR: ForgotPassword failed for email %s: %v", req.Email, err)
	}

	ac.JSONSuccess(c, gin.H{"message": "If that email exists, a reset link has been sent."})
}

func (ac *LocalAuthController) ResetPassword(c *gin.Context) {
	var req auth.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.JSONValidationError(c, "Invalid request payload")
		return
	}

	ctx := context.Background()
	err := ac.authService.ResetPassword(ctx, req)
	if err != nil {
		ac.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "Failed to reset password. Token may be invalid.")
		return
	}

	ac.JSONSuccess(c, gin.H{"message": "Password reset successful."})
}

func (ac *LocalAuthController) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		ac.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	var req auth.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.JSONValidationError(c, "Invalid request payload")
		return
	}

	ctx := context.Background()
	err := ac.authService.ChangePassword(ctx, userID.(string), req)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			ac.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid old password")
			return
		}
		ac.JSONInternalError(c, "Failed to change password")
		return
	}

	ac.JSONSuccess(c, gin.H{"message": "Password changed successfully."})
}

func (ac *LocalAuthController) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	ac.JSONSuccess(c, gin.H{"message": "Logged out successfully"})
}
