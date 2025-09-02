package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/devfajar/task-management-system/internal/delivery/http/middleware"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	UC   usecases.TaskUsecase
	Auth middleware.Auth
}

func (h *TaskHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	api := v1.Group("/tasks")
	{
		api.POST("", h.Auth.Require(helper.TasksCreate), h.create)
		api.PUT("/:id", h.Auth.Require(helper.TasksUpdate), h.updateFields)
		api.PATCH("/:id/estimates", h.Auth.Require(helper.TasksUpdate), h.updateEstimates)
		api.PUT("/:id/move", h.Auth.Require(helper.TasksUpdate), h.move)
		api.PUT("/:id/reorder", h.Auth.Require(helper.TasksUpdate), h.reorder)
	}

	// list tasks by column
	v1.GET("/columns/:column_id/tasks", h.listByColumn)
}

type taskCreateReq struct {
	ProjectID                string   `json:"project_id" binding:"required,uuid"`
	BoardID                  string   `json:"board_id" binding:"required,uuid"`
	ColumnID                 string   `json:"column_id" binding:"required,uuid"`
	Title                    string   `json:"title" binding:"required"`
	Description              *string  `json:"description"`
	Status                   *string  `json:"status"`
	Priority                 *int32   `json:"priority"`
	AssigneeID               *string  `json:"assignee_id"`
	DueDate                  *string  `json:"due_date"` // RFC3339
	Position                 *float64 `json:"position"`
	StoryPoints              *float64 `json:"story_points"`
	OriginalEstimateSeconds  *int32   `json:"original_estimate_seconds"`
	RemainingEstimateSeconds *int32   `json:"remaining_estimate_seconds"`
}

func (h *TaskHandler) create(c *gin.Context) {
	var req taskCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pid, _ := uuid.Parse(req.ProjectID)
	bid, _ := uuid.Parse(req.BoardID)
	cid, _ := uuid.Parse(req.ColumnID)

	var ass *uuid.UUID
	if req.AssigneeID != nil {
		a, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id"})
			return
		}
		ass = &a
	}
	var due *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date (RFC3339)"})
			return
		}
		due = &t
	}

	task, err := h.UC.Create(c, repository.CreateTaskParams{
		ProjectID:                pid,
		BoardID:                  bid,
		ColumnID:                 cid,
		Title:                    req.Title,
		Description:              req.Description,
		Status:                   req.Status,
		Priority:                 req.Priority,
		AssigneeID:               ass,
		DueDate:                  due,
		Position:                 req.Position,
		StoryPoints:              req.StoryPoints,
		OriginalEstimateSeconds:  req.OriginalEstimateSeconds,
		RemainingEstimateSeconds: req.RemainingEstimateSeconds,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

type taskUpdateReq struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Status      *string  `json:"status"`
	Priority    *int32   `json:"priority"`
	AssigneeID  *string  `json:"assignee_id"`
	DueDate     *string  `json:"due_date"`
	StoryPoints *float64 `json:"story_points"`
}

func (h *TaskHandler) updateFields(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req taskUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ass *uuid.UUID
	if req.AssigneeID != nil {
		a, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id"})
			return
		}
		ass = &a
	}
	var due *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date"})
			return
		}
		due = &t
	}

	task, err := h.UC.UpdateFields(c, repository.UpdateTaskFieldsParams{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssigneeID:  ass,
		DueDate:     due,
		StoryPoints: req.StoryPoints,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

type estReq struct {
	OriginalEstimateSeconds  *int32 `json:"original_estimate_seconds"`
	RemainingEstimateSeconds *int32 `json:"remaining_estimate_seconds"`
}

func (h *TaskHandler) updateEstimates(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req estReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.UC.UpdateEstimates(c, id, req.OriginalEstimateSeconds, req.RemainingEstimateSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

type moveReq struct {
	ToColumnID string   `json:"to_column_id" binding:"required,uuid"`
	NewPos     *float64 `json:"new_position"`
}

func (h *TaskHandler) move(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req moveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	col, _ := uuid.Parse(req.ToColumnID)
	if err := h.UC.Move(c, id, col, req.NewPos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

type reorderReq struct {
	LeftPos  *float64 `json:"left_pos"`
	RightPos *float64 `json:"right_pos"`
}

func (h *TaskHandler) reorder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req reorderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.UC.Reorder(c, id, req.LeftPos, req.RightPos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) listByColumn(c *gin.Context) {
	cid, err := uuid.Parse(c.Param("column_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid column_id"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.UC.ListByColumn(c, cid, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
