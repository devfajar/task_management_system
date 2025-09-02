package handler

import (
	"net/http"

	"github.com/devfajar/task-management-system/internal/delivery/http/middleware"
	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ColumnHandler struct {
	UC   usecases.ColumnUsecase
	Auth middleware.Auth
}

func (h *ColumnHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	// nested di board
	v1.POST("/boards/:board_id/columns", h.Auth.Require(helper.ProjectsUpdate), h.create)
	v1.GET("/boards/:board_id/columns", h.list)
	// standalone by id
	v1.PUT("/columns/:id", h.Auth.Require(helper.ProjectsUpdate), h.update)
	v1.DELETE("/columns/:id", h.Auth.Require(helper.ProjectsUpdate), h.softDelete)
}

type colCreateReq struct {
	Name     string   `json:"name" binding:"required"`
	WIPLimit *int32   `json:"wip_limit"`
	Position *float64 `json:"position"`
}

func (h *ColumnHandler) create(c *gin.Context) {
	bid, err := uuid.Parse(c.Param("board_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board_id"})
		return
	}
	var req colCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	col, err := h.UC.Create(c, bid, req.Name, req.WIPLimit, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, col)
}

func (h *ColumnHandler) list(c *gin.Context) {
	bid, err := uuid.Parse(c.Param("board_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board_id"})
		return
	}
	cl, err := h.UC.ListByBoard(c, bid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cl)
}

type colUpdateReq struct {
	Name     *string  `json:"name"`
	WIPLimit *int32   `json:"wip_limit"`
	Position *float64 `json:"position"`
}

func (h *ColumnHandler) update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req colUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	col, err := h.UC.Update(c, id, req.Name, req.WIPLimit, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, col)
}

func (h *ColumnHandler) softDelete(c *gin.Context) {
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
