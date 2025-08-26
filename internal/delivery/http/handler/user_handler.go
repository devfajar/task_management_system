package handler

import (
	"net/http"
	"strconv"

	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct{ UC usecases.UserUsecase }

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/v1/users")
	{
		api.POST("", h.create)
		api.GET("", h.list)
		api.GET("/:id", h.get)
		api.PUT("/:id", h.updateProfile)
		api.PATCH("/:id/password", h.updatePassword)
		api.DELETE("/:id", h.delete)
		api.POST("/:id/restore", h.restore)
	}
}

// TODO Must Use DTO
type createReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type pwReq struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type updateReq struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	IsActive *bool   `json:"is_active,omitempty"`
}

func (h *UserHandler) create(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.UC.Create(c, req.Username, req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case usecases.ErrBadInput:
			status = http.StatusBadRequest
		case usecases.ErrEmailExists:
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	u, err := h.UC.GetById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) list(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	us, err := h.UC.List(c, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, us)
}

func (h *UserHandler) updateProfile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.UC.UpdateProfile(c, id, req.Username, req.Email, req.IsActive)
	if err != nil {
		status := http.StatusInternalServerError
		if err == usecases.ErrEmailExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) updatePassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req pwReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.UC.UpdatePassword(c, id, req.NewPassword); err != nil {
		status := http.StatusInternalServerError
		if err == usecases.ErrBadInput {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	force := c.Query("force") == "true"
	var opErr error
	if force {
		// TODO: pastikan hanya admin yang boleh force delete (middleware/authorization)
		opErr = h.UC.ForceDelete(c, id)
	} else {
		opErr = h.UC.Delete(c, id)
	}
	if opErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": opErr.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) restore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.UC.Restore(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
