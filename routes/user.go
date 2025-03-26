// routes/user.go
package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/shakyasimha/go-backend-server/models"
	"github.com/shakyasimha/go-backend-server/utils"

	"golang.org/x/crypto/bcrypt"
)

func UserRoutes(router *gin.Engine) {
	user := models.NewUser()
	db := user.ConnectDB()

	route := router.Group("/users")

	// signup is for creating a new user
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

		token, err := utils.GenerateToken(user.ID)

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(201, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"token":    token,
		})
	})

	// Logging in the user
	route.POST("/login", func(c *gin.Context) {
		var input models.Input

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		var user models.User

		if result := db.Where("username = ?", input.Username).First(&user); result.Error != nil {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		token, err := utils.GenerateToken(user.ID)

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, gin.H{
			"message":  "Login successful",
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"token":    token,
		})
	})
}
