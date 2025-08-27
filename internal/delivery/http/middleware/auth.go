package middleware

import (
	"net/http"

	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Auth struct {
	Perms usecases.PermissionUsecase
}

// DevAuth: ambil user dari header X-User-ID (untuk dev). Production: ganti JWT.
func (auth Auth) DevAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-User-ID")
		if id == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID"})
			return
		}
		if _, err := uuid.Parse(id); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid X-User-ID"})
			return
		}
		c.Set("user_id", id)
		c.Next()
	}
}

func (auth Auth) Require(permKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth"})
			return
		}
		uid, _ := uuid.Parse(val.(string))
		okPerm, err := auth.Perms.Has(c, uid, permKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !okPerm {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: " + permKey})
			return
		}
		c.Next()
	}
}
