package handlers

import (
	"github.com/devchuckcamp/goauthx"
	"github.com/devchuckcamp/gocommerce-api/internal/http/middleware"
	"github.com/devchuckcamp/gocommerce-api/internal/http/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *goauthx.Service
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *goauthx.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req goauthx.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	authResp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if err == goauthx.ErrEmailAlreadyExists {
			response.Conflict(c, "Email already exists")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, gin.H{
		"user":          authResp.User,
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
		"expires_at":    authResp.ExpiresAt,
	})
}

// Login handles user login
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req goauthx.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	authResp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if err == goauthx.ErrInvalidCredentials || err == goauthx.ErrUserNotFound {
			response.Unauthorized(c, "Invalid credentials")
			return
		}
		if err == goauthx.ErrUserInactive {
			response.Forbidden(c, "Account is inactive")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"user":          authResp.User,
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
		"expires_at":    authResp.ExpiresAt,
	})
}

// Profile handles retrieving user profile
// GET /auth/profile
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	roles, _ := middleware.GetUserRoles(c)

	response.Success(c, gin.H{
		"user":  user,
		"roles": roles,
	})
}

// RefreshToken handles token refresh
// POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req goauthx.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	authResp, err := h.authService.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == goauthx.ErrInvalidRefreshToken {
			response.Unauthorized(c, "Invalid refresh token")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
		"expires_at":    authResp.ExpiresAt,
	})
}

// Logout handles user logout
// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.NoContent(c)
}

// GoogleOAuthURL generates the Google OAuth authorization URL
// GET /auth/google
func (h *AuthHandler) GoogleOAuthURL(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		// Generate a random state for CSRF protection
		state = "random-state-" + c.Request.RemoteAddr
	}

	url, err := h.authService.GetGoogleOAuthURL(goauthx.GoogleOAuthURLRequest{
		State: state,
	})
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"url":   url,
		"state": state,
	})
}

// GoogleOAuthCallback handles the Google OAuth callback
// GET /auth/google/callback
func (h *AuthHandler) GoogleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		errorDesc := c.Query("error_description")
		response.BadRequest(c, "OAuth error: "+errorParam+" - "+errorDesc)
		return
	}

	if code == "" {
		response.BadRequest(c, "Missing authorization code")
		return
	}

	authResp, err := h.authService.HandleGoogleOAuthCallback(c.Request.Context(), goauthx.GoogleOAuthCallbackRequest{
		Code:  code,
		State: state,
	})
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"user":          authResp.User,
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
		"expires_at":    authResp.ExpiresAt,
	})
}
