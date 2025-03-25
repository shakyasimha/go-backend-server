package utils

import (
	"log"
	"time"

	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/models"
	"github.com/shakyasimha/go-backend-server/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		db := user.ConnectDB()

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Basic ") {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		creds, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": "Invalid Base64 encoding"})
			c.Abort()
			return
		}

		username, password := utils.SplitCredentials(string(creds))
		if username == "" || password == "" {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": "Invalid credentials format"})
			c.Abort()
			return
		}

		if result := db.Where("username = ?", username).First(&user); result.Error != nil {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		if user.Password != password { // Insecure; use bcrypt in production
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		// Set user_id in context for use in routes
		c.Set("user_id", user.ID)
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
