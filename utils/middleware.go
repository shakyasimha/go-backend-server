package utils

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apikey := c.GetHeader("X-API-KEY")
		if apikey != "secret-key" { // Example: require a specific key
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or missing API key"})
			return
		}
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request - Method: %s | Status: %d | Duration: %v", c.Request.Method, c.Writer.Status(), duration)
	}
}
