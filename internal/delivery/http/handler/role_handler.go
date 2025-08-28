package handler

import (
	"net/http"

	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	UC usecases.RoleUsecase
}

func (h *RoleHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	api := v1.Group("/roles")
	{
		api.POST("", h.create)
		api.GET("", h.list)
		api.DELETE("/:id", h.delete)

		api.POST("/assign", h.assign) // body: { "user_id": "...", "role_key": "..." }
		api.POST("/revoke", h.revoke) // body: { "user_id": "...", "role_key": "..." }
	}
}

func (h *RoleHandler) create(c *gin.Context) {
	var req struct{ Key, Name string }
	if err := c.ShouldBindJSON(&req); err != nil || req.Key == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key & name required"})
		return
	}
	role, err := h.UC.CreateRole(c, req.Key, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, role)
}

func (h *RoleHandler) list(c *gin.Context) {
	roles, err := h.UC.ListRoles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.UC.DeleteRole(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RoleHandler) assign(c *gin.Context) {
	var req struct {
		UserID  string `json:"user_id" binding:"required"`
		RoleKey string `json:"role_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	if err := h.UC.Assign(c, uid, req.RoleKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RoleHandler) revoke(c *gin.Context) {
	var req struct {
		UserID  string `json:"user_id" binding:"required"`
		RoleKey string `json:"role_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	if err := h.UC.Revoke(c, uid, req.RoleKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
