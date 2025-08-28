package handler

import (
	"net/http"

	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct{ UC usecases.PermissionUsecase }

func (h *PermissionHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	api := v1.Group("/permissions")
	{
		api.POST("", h.create)
		api.GET("", h.list)
		api.DELETE("/:id", h.delete)

		api.POST("/grant", h.grant)
		api.POST("/revoke", h.revoke)
	}
}

func (h *PermissionHandler) create(c *gin.Context) {
	var req struct{ Key, Description string }
	if err := c.ShouldBindJSON(&req); err != nil || req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key  required"})
	}
	perm, err := h.UC.CreatePermission(c, req.Key, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, perm)
}

func (h *PermissionHandler) list(c *gin.Context) {
	perms, err := h.UC.ListPermissions(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perms)
}

func (h *PermissionHandler) delete(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (h *PermissionHandler) grant(c *gin.Context) {
	var req struct {
		RoleKey string `json:"role_key" binding:"required"`
		PermKey string `json:"perm_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RoleKey == "" || req.PermKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_key & perm_key required"})
		return
	}
	if err := h.UC.Grant(c, req.RoleKey, req.PermKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PermissionHandler) revoke(c *gin.Context) {
	var req struct {
		RoleKey string `json:"role_key" binding:"required"`
		PermKey string `json:"perm_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RoleKey == "" || req.PermKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_key & perm_key required"})
		return
	}
	if err := h.UC.Revoke(c, req.RoleKey, req.PermKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
