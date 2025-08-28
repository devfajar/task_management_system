package middleware

import (
	"net/http"
	"strings"

	"github.com/devfajar/task-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JWTConfig struct {
	Secret []byte
}

func JWT(cfg JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tok := strings.TrimSpace(h[7:])
		claims, err := utils.ParseAndVerify(tok, cfg.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		if claims.UserID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid subject"})
			return
		}
		c.Set("user_id", claims.UserID.String())
		c.Set("roles", claims.Roles)
		c.Set("perms", claims.Perms)
		c.Next()
	}
}
