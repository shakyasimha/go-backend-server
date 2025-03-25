package utils

import (
	"errors"
	"log"
	"time"

	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/shakyasimha/go-backend-server/models"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		username, password, err := SplitCredentials(string(creds))
		if err != nil {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		var user models.User
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

		// Set user_id as uint (from gorm.Model.ID)
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

// Helper function to split username:password
func SplitCredentials(creds string) (username, password string, err error) {
	if creds == "" {
		return "", "", errors.New("credentials string is empty")
	}
	parts := strings.SplitN(creds, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid credentials format: expected 'username:password'")
	}
	return parts[0], parts[1], nil
}
