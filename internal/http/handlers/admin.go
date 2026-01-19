package handlers

import (
	"github.com/devchuckcamp/goauthx"
	"github.com/devchuckcamp/goauthx/pkg/rbac"
	"github.com/devchuckcamp/gocommerce-api/internal/http/response"
	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin RBAC management endpoints
type AdminHandler struct {
	authService *goauthx.Service
	authStore   goauthx.Store
	seeder      *goauthx.Seeder
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(authService *goauthx.Service, authStore goauthx.Store, seeder *goauthx.Seeder) *AdminHandler {
	return &AdminHandler{
		authService: authService,
		authStore:   authStore,
		seeder:      seeder,
	}
}

// --- Role Management ---

// ListRoles returns all roles
// GET /admin/roles
func (h *AdminHandler) ListRoles(c *gin.Context) {
	roles, err := h.authStore.ListRoles(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to list roles")
		return
	}

	response.Success(c, gin.H{"roles": roles})
}

// CreateRole creates a new role
// POST /admin/roles
func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	role := &goauthx.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.authStore.CreateRole(c.Request.Context(), role); err != nil {
		response.InternalServerError(c, "Failed to create role")
		return
	}

	response.Created(c, gin.H{"role": role})
}

// GetRole returns a specific role by ID
// GET /admin/roles/:id
func (h *AdminHandler) GetRole(c *gin.Context) {
	roleID := c.Param("id")

	role, err := h.authStore.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		response.NotFound(c, "Role not found")
		return
	}

	response.Success(c, gin.H{"role": role})
}

// UpdateRole updates an existing role
// PUT /admin/roles/:id
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	role, err := h.authStore.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		response.NotFound(c, "Role not found")
		return
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := h.authStore.UpdateRole(c.Request.Context(), role); err != nil {
		response.InternalServerError(c, "Failed to update role")
		return
	}

	response.Success(c, gin.H{"role": role})
}

// DeleteRole deletes a role
// DELETE /admin/roles/:id
func (h *AdminHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")

	if err := h.authStore.DeleteRole(c.Request.Context(), roleID); err != nil {
		response.InternalServerError(c, "Failed to delete role")
		return
	}

	response.NoContent(c)
}

// --- Permission Management ---

// ListPermissions returns all permissions
// GET /admin/permissions
func (h *AdminHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.authStore.ListPermissions(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to list permissions")
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}

// CreatePermission creates a new permission
// POST /admin/permissions
func (h *AdminHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Resource    string `json:"resource" binding:"required"`
		Action      string `json:"action" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	permission := &goauthx.Permission{
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	if err := h.authStore.CreatePermission(c.Request.Context(), permission); err != nil {
		response.InternalServerError(c, "Failed to create permission")
		return
	}

	response.Created(c, gin.H{"permission": permission})
}

// GetPermission returns a specific permission by ID
// GET /admin/permissions/:id
func (h *AdminHandler) GetPermission(c *gin.Context) {
	permissionID := c.Param("id")

	permission, err := h.authStore.GetPermissionByID(c.Request.Context(), permissionID)
	if err != nil {
		response.NotFound(c, "Permission not found")
		return
	}

	response.Success(c, gin.H{"permission": permission})
}

// UpdatePermission updates an existing permission
// PUT /admin/permissions/:id
func (h *AdminHandler) UpdatePermission(c *gin.Context) {
	permissionID := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Resource    string `json:"resource"`
		Action      string `json:"action"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	permission, err := h.authStore.GetPermissionByID(c.Request.Context(), permissionID)
	if err != nil {
		response.NotFound(c, "Permission not found")
		return
	}

	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Resource != "" {
		permission.Resource = req.Resource
	}
	if req.Action != "" {
		permission.Action = req.Action
	}
	if req.Description != "" {
		permission.Description = req.Description
	}

	if err := h.authStore.UpdatePermission(c.Request.Context(), permission); err != nil {
		response.InternalServerError(c, "Failed to update permission")
		return
	}

	response.Success(c, gin.H{"permission": permission})
}

// DeletePermission deletes a permission
// DELETE /admin/permissions/:id
func (h *AdminHandler) DeletePermission(c *gin.Context) {
	permissionID := c.Param("id")

	if err := h.authStore.DeletePermission(c.Request.Context(), permissionID); err != nil {
		response.InternalServerError(c, "Failed to delete permission")
		return
	}

	response.NoContent(c)
}

// --- User Role Assignments ---

// GetUserRoles returns all roles assigned to a user
// GET /admin/users/:id/roles
func (h *AdminHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("id")

	roles, err := h.authService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Failed to get user roles")
		return
	}

	response.Success(c, gin.H{"roles": roles})
}

// AssignRoleToUser assigns a role to a user
// POST /admin/users/:id/roles
func (h *AdminHandler) AssignRoleToUser(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		RoleID   string `json:"role_id"`
		RoleName string `json:"role_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if req.RoleID == "" && req.RoleName == "" {
		response.BadRequest(c, "Either role_id or role_name is required")
		return
	}

	var roleName string
	if req.RoleName != "" {
		roleName = req.RoleName
	} else {
		role, err := h.authStore.GetRoleByID(c.Request.Context(), req.RoleID)
		if err != nil {
			response.NotFound(c, "Role not found")
			return
		}
		roleName = role.Name
	}

	if err := h.seeder.AssignRoleToUser(c.Request.Context(), userID, rbac.RoleName(roleName)); err != nil {
		response.InternalServerError(c, "Failed to assign role to user")
		return
	}

	response.Success(c, gin.H{"message": "Role assigned successfully"})
}

// RemoveRoleFromUser removes a role from a user
// DELETE /admin/users/:id/roles/:roleId
func (h *AdminHandler) RemoveRoleFromUser(c *gin.Context) {
	userID := c.Param("id")
	roleID := c.Param("roleId")

	// Get role name from role ID
	role, err := h.authStore.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		response.NotFound(c, "Role not found")
		return
	}

	if err := h.seeder.RemoveRoleFromUser(c.Request.Context(), userID, rbac.RoleName(role.Name)); err != nil {
		response.InternalServerError(c, "Failed to remove role from user")
		return
	}

	response.NoContent(c)
}

// --- Role Permission Grants ---

// GetRolePermissions returns all permissions granted to a role
// GET /admin/roles/:id/permissions
func (h *AdminHandler) GetRolePermissions(c *gin.Context) {
	roleID := c.Param("id")

	permissions, err := h.authStore.GetRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		response.InternalServerError(c, "Failed to get role permissions")
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}

// GrantPermissionToRole grants a permission to a role
// POST /admin/roles/:id/permissions
func (h *AdminHandler) GrantPermissionToRole(c *gin.Context) {
	roleID := c.Param("id")

	var req struct {
		PermissionID string `json:"permission_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.authStore.GrantPermission(c.Request.Context(), roleID, req.PermissionID); err != nil {
		response.InternalServerError(c, "Failed to grant permission to role")
		return
	}

	response.Success(c, gin.H{"message": "Permission granted successfully"})
}

// RevokePermissionFromRole revokes a permission from a role
// DELETE /admin/roles/:id/permissions/:permId
func (h *AdminHandler) RevokePermissionFromRole(c *gin.Context) {
	roleID := c.Param("id")
	permissionID := c.Param("permId")

	if err := h.authStore.RevokePermission(c.Request.Context(), roleID, permissionID); err != nil {
		response.InternalServerError(c, "Failed to revoke permission from role")
		return
	}

	response.NoContent(c)
}
