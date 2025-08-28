package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
}

func (Auth) Require(permKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("perms")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing perms on token"})
			return
		}
		perms, _ := val.([]string)
		for _, p := range perms {
			if p == permKey {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: " + permKey})
	}
}
