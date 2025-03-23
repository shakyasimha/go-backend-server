package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apikey := c.GetHeader("X-API-KEY")
		if apikey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "API key is missing"})
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

func main() {
	// Create a default Gin router
	router := gin.Default()

	// Define a route for the root URL
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	// AuthMiddleware for specific routes
	authGroup := router.Group("/api")
	authGroup.Use(AuthMiddleware())
	{
		authGroup.GET("/data", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Authenticated and authorized!"})
		})
	}

	// Run the server on port 8080
	router.Run(":8080")
}
