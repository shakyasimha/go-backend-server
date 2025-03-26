// utils/middleware.go
package utils

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var jwtSecret = []byte("your-secret-key") // Use env variable in production

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token: " + err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			userIDFloat, ok := (*claims)["user_id"].(float64) // JWT stores numbers as float64

			if !ok {
				c.JSON(401, gin.H{"error": "Invalid user_id in token"})
				c.Abort()
				return
			}

			userID := uint(userIDFloat)
			c.Set("user_id", userID)
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("Request - Method: %s | Status: %d | Duration: %v",
			c.Request.Method, c.Writer.Status(), duration)
	}
}

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
