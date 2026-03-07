package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panicked: %v\n%s", err, debug.Stack())
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
