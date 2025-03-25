package routes

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/models"
	"github.com/shakyasimha/go-backend-server/utils"
)

func UserRoutes(router *gin.Engine) {
	user := models.NewUser()
	db := user.ConnectDB()

	route := router.Group("/users")

	// POST /users/signup - unchanged
	route.POST("/signup", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}
		if results := db.Create(&user); results.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create user: " + results.Error.Error()})
			return
		}
		c.JSON(201, gin.H{"id": user.ID, "username": user.Username, "email": user.Email})
	})

	// POST /users/login
	route.POST("/login", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Basic ") {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		creds, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid Base64 encoding"})
			return
		}

		// Fix: Handle all 3 return values
		username, password, err := utils.SplitCredentials(string(creds))
		if err != nil || username == "" || password == "" {
			c.JSON(400, gin.H{"error": "Invalid credentials format"})
			return
		}

		var user models.User
		if result := db.Where("username = ?", username).First(&user); result.Error != nil {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		if user.Password != password { // Insecure; use bcrypt in production
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		c.JSON(200, gin.H{
			"message":  "Login successful",
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		})
	})
}
