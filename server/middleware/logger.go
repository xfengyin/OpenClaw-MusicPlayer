package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("[%s] %s %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())

		c.Next()

		log.Printf("[%s] Completed in %v", c.Request.Method, time.Since(start))
	}
}
