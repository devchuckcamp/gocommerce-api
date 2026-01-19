package middleware

import (
	"strings"

	"github.com/devchuckcamp/goauthx"
	"github.com/devchuckcamp/gocommerce-api/internal/http/response"
	"github.com/gin-gonic/gin"
)

const (
	// UserIDKey is the context key for user ID
	UserIDKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey = "user_email"
	// UserRolesKey is the context key for user roles
	UserRolesKey = "user_roles"
)

// AuthMiddleware wraps goauthx authentication for Gin
type AuthMiddleware struct {
	authService *goauthx.Service
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authService *goauthx.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate validates JWT tokens and sets user context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token using goauthx
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserRolesKey, claims.Roles)

		c.Next()
	}
}

// RequireRole checks if the authenticated user has a specific role
func (m *AuthMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get(UserRolesKey)
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			response.InternalServerError(c, "Invalid roles in context")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range userRoles {
			if role == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole checks if the authenticated user has any of the specified roles
func (m *AuthMiddleware) RequireAnyRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get(UserRolesKey)
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			response.InternalServerError(c, "Invalid roles in context")
			c.Abort()
			return
		}

		hasRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range roleNames {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID extracts the user ID from the Gin context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUserEmail extracts the user email from the Gin context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

// GetUserRoles extracts the user roles from the Gin context
func GetUserRoles(c *gin.Context) ([]string, bool) {
	roles, exists := c.Get(UserRolesKey)
	if !exists {
		return nil, false
	}
	r, ok := roles.([]string)
	return r, ok
}

// UserIDExtractor is a function type that extracts a resource owner's user ID from the request
type UserIDExtractor func(c *gin.Context) (string, error)

// RequirePermission checks if the authenticated user has a specific permission
func (m *AuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		hasPermission, err := m.authService.HasPermission(c.Request.Context(), userID, permission)
		if err != nil {
			response.InternalServerError(c, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if the authenticated user has any of the specified permissions
func (m *AuthMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		hasAny, err := m.authService.HasAnyPermission(c.Request.Context(), userID, permissions)
		if err != nil {
			response.InternalServerError(c, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasAny {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions checks if the authenticated user has all of the specified permissions
func (m *AuthMiddleware) RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		hasAll, err := m.authService.HasAllPermissions(c.Request.Context(), userID, permissions)
		if err != nil {
			response.InternalServerError(c, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasAll {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOwnerOrRole checks if the user owns the resource OR has one of the specified roles
func (m *AuthMiddleware) RequireOwnerOrRole(getOwnerID UserIDExtractor, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// Check if user is the owner
		ownerID, err := getOwnerID(c)
		if err == nil && ownerID == userID {
			c.Next()
			return
		}

		// If not owner, check if user has any of the required roles
		userRoles, ok := GetUserRoles(c)
		if !ok {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		for _, userRole := range userRoles {
			for _, requiredRole := range roles {
				if userRole == requiredRole {
					c.Next()
					return
				}
			}
		}

		response.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}

// RequireOwnerOrPermission checks if the user owns the resource OR has the specified permission
func (m *AuthMiddleware) RequireOwnerOrPermission(getOwnerID UserIDExtractor, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// Check if user is the owner
		ownerID, err := getOwnerID(c)
		if err == nil && ownerID == userID {
			c.Next()
			return
		}

		// If not owner, check if user has the required permission
		hasPermission, err := m.authService.HasPermission(c.Request.Context(), userID, permission)
		if err != nil {
			response.InternalServerError(c, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
