package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/RubenPari/clear-songs/internal/application/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthControllerRefactored is the refactored auth controller using dependency injection
type AuthControllerRefactored struct {
	BaseController
	loginUC    *auth.LoginUseCase
	callbackUC *auth.CallbackUseCase
	logoutUC   *auth.LogoutUseCase
	isAuthUC   *auth.IsAuthUseCase
}

// NewAuthControllerRefactored creates a new auth controller
func NewAuthControllerRefactored(
	loginUC *auth.LoginUseCase,
	callbackUC *auth.CallbackUseCase,
	logoutUC *auth.LogoutUseCase,
	isAuthUC *auth.IsAuthUseCase,
) *AuthControllerRefactored {
	return &AuthControllerRefactored{
		loginUC:    loginUC,
		callbackUC: callbackUC,
		logoutUC:   logoutUC,
		isAuthUC:   isAuthUC,
	}
}

// Login handles GET /auth/login
func (ac *AuthControllerRefactored) Login(c *gin.Context) {
	url := ac.loginUC.Execute()
	c.Redirect(302, url)
}

// Callback handles GET /auth/callback
func (ac *AuthControllerRefactored) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		ac.JSONValidationError(c, "Authorization code is required")
		return
	}

	// Try to get local User ID from JWT if present
	localUserID := ""
	tokenString, errCookie := c.Cookie("auth_token")
	if errCookie == nil && tokenString != "" {
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "fallback-secret-for-dev"
		}

		token, errJWT := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if errJWT == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if sub, ok := claims["sub"].(string); ok {
					localUserID = sub
				}
			}
		}
	}

	ctx := context.Background()
	redirectURL, err := ac.callbackUC.Execute(ctx, code, localUserID)
	if err != nil {
		ac.JSONInternalError(c, "Error authenticating user")
		return
	}

	c.Redirect(302, redirectURL)
}

// Logout handles GET /auth/logout
func (ac *AuthControllerRefactored) Logout(c *gin.Context) {
	ctx := context.Background()
	if err := ac.logoutUC.Execute(ctx); err != nil {
		ac.JSONInternalError(c, "Error logging out")
		return
	}

	ac.JSONSuccess(c, gin.H{"message": "User logged out successfully"})
}

// IsAuth handles GET /auth/is-auth
func (ac *AuthControllerRefactored) IsAuth(c *gin.Context) {
	ctx := context.Background()
	userInfo, err := ac.isAuthUC.Execute(ctx)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	ac.JSONSuccess(c, gin.H{
		"user": gin.H{
			"spotify_id":    userInfo.SpotifyID,
			"display_name":  userInfo.DisplayName,
			"email":         userInfo.Email,
			"profile_image": userInfo.ProfileImage,
		},
	})
}
