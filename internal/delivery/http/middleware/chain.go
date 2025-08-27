package middleware

import "github.com/gin-gonic/gin"

func Chain(mw gin.HandlerFunc) func(gin.HandlerFunc) gin.HandlerFunc {
	return func(h gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			mw(c)
			if !c.IsAborted() {
				h(c)
			}
		}
	}
}
