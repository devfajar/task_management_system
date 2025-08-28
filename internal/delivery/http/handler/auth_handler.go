package handler

import (
	"net/http"

	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UC usecases.AuthUsecase
}

func (h *AuthHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	api := v1.Group("/auth")
	{
		api.POST("/login", h.login)
		api.POST("/refresh", h.refresh)
		api.POST("/logout", h.logout)
	}
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ua := c.GetHeader("User-Agent")
	ip := c.ClientIP()
	access, refresh, err := h.UC.Login(c, req.Email, req.Password, ua, ip)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"token_type":    "Bearer",
	})
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ua := c.GetHeader("User-Agent")
	ip := c.ClientIP()
	access, refresh, err := h.UC.Refresh(c, req.RefreshToken, ua, ip)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"token_type":    "Bearer",
	})
}

type logoutReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) logout(c *gin.Context) {
	var req logoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.UC.Logout(c, req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
