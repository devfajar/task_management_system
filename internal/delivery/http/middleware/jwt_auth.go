package middleware

import (
	"net/http"
	"strings"

	"github.com/devfajar/task-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JWTConfig struct {
	Secret            []byte
	CheckTokenVersion bool
	GetTokenVersion   func(userID uuid.UUID) (int, error)
}

func JWT(cfg JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimSpace(h[7:])
		claims, err := utils.ParseAndVerify(token, cfg.Secret)
		if err != nil || claims.UserID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		if cfg.CheckTokenVersion && cfg.GetTokenVersion != nil {
			tv, err := cfg.GetTokenVersion(claims.UserID)
			if err != nil || int32(tv) != int32(claims.TokenVersion) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
				return
			}
		}
		c.Set("user_id", claims.UserID.String())
		c.Set("roles", claims.Roles)
		c.Set("perms", claims.Perms)
		c.Next()
	}
}
