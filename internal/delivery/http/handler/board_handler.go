package handler

import (
	"net/http"

	"github.com/devfajar/task-management-system/internal/delivery/http/middleware"
	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BoardHandler struct {
	UC   usecases.BoardUsecase
	Auth middleware.Auth
}

func (h *BoardHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	api := v1.Group("/boards")
	{
		api.POST("", h.Auth.Require(helper.BoardsCreate), h.create)
		api.GET("", h.listByProject)
		api.GET("/:board_id", h.get)
		api.DELETE("/:board_id", h.Auth.Require(helper.BoardsDelete), h.softDelete)
		api.DELETE("/:board_id/force", h.Auth.Require(helper.BoardsDelete), h.forceDelete)
	}
}

type boardCreateReq struct {
	ProjectID   string  `json:"project_id" binding:"required,uuid"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

func (h *BoardHandler) create(c *gin.Context) {
	var req boardCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pid, _ := uuid.Parse(req.ProjectID)
	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}
	b, err := h.UC.Create(c, pid, req.Name, desc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *BoardHandler) listByProject(c *gin.Context) {
	pidStr := c.Query("project_id")
	if pidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id required"})
		return
	}
	pid, err := uuid.Parse(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}
	list, err := h.UC.ListByProject(c, pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *BoardHandler) get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("board_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	b, err := h.UC.Get(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BoardHandler) softDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("board_id"))
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

func (h *BoardHandler) forceDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("board_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.UC.ForceDelete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
