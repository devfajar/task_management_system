package handler

import (
	"net/http"
	"strconv"

	"github.com/devfajar/task-management-system/internal/delivery/http/middleware"
	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	UC    usecases.ProjectUsecase
	Auth  middleware.Auth
	Audit usecases.AuditUsecase // optional: log force delete
}

func (h *ProjectHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	api := v1.Group("/projects")
	{
		api.POST("", h.Auth.Require(helper.ProjectsCreate), h.create)
		api.GET("", h.list)
		api.GET("/:id", h.get)
		api.PUT("/:id", h.Auth.Require(helper.ProjectsUpdate), h.update)
		api.DELETE("/:id", h.Auth.Require(helper.ProjectsDelete), h.softDelete)
		api.DELETE("/:id/force", h.Auth.Require(helper.ProjectsDelete), h.forceDelete)
	}
}

type projectCreateReq struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	OwnerID     *string `json:"owner_id"`
}

type projectUpdateReq struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (h *ProjectHandler) create(c *gin.Context) {
	var req projectCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var owner *uuid.UUID
	if req.OwnerID != nil && *req.OwnerID != "" {
		id, err := uuid.Parse(*req.OwnerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid owner_id"})
			return
		}
		owner = &id
	} else {
		// default: owner = current user
		if uidStr, ok := c.Get("user_id"); ok {
			if id, err := uuid.Parse(uidStr.(string)); err == nil {
				owner = &id
			}
		}
	}
	p, err := h.UC.Create(c, req.Name, req.Description, owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProjectHandler) list(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.UC.List(c, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ProjectHandler) get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.UC.Get(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req projectUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.UC.Update(c, id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) softDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.UC.SoftDelete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) forceDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.UC.ForceDelete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// audit optional
	if h.Audit != nil {
		if uidStr, ok := c.Get("user_id"); ok {
			if uid, e := uuid.Parse(uidStr.(string)); e == nil {
				_ = h.Audit.Log(c, &uid, "projects.force_delete", "projects", &id, nil, c.ClientIP(), c.GetHeader("User-Agent"))
			}
		}
	}
	c.Status(http.StatusNoContent)
}
