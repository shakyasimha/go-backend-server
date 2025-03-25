package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/models"
	"github.com/shakyasimha/go-backend-server/utils"
	"gorm.io/gorm"
)

func TodoRoutes(router *gin.Engine) {
	todo := models.NewTodo()
	db := todo.ConnectDB()

	route := router.Group("/todos", utils.AuthMiddleware(db))

	// Create
	route.POST("/", func(c *gin.Context) {
		var todo models.Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(500, gin.H{"error": "User ID not found in context"})
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			c.JSON(500, gin.H{"error": "User ID is not a uint"})
			return
		}

		// Fix: Convert uint to string
		todo.UserID = strconv.FormatUint(uint64(uid), 10)

		if result := db.Create(&todo); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create todo: " + result.Error.Error()})
			return
		}
		c.JSON(200, todo)
	})

	// List (only user's todos)
	route.GET("/", func(c *gin.Context) {
		var todos []models.Todo
		userID, exists := c.Get("user_id")

		if !exists {
			c.JSON(500, gin.H{"error": "User ID not found in context"})
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(500, gin.H{"error": "User ID is not a uint"})
			return
		}
		// Convert uid to string for query
		if result := db.Where("user_id = ?", strconv.FormatUint(uint64(uid), 10)).Find(&todos); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve todos: " + result.Error.Error()})
			return
		}
		c.JSON(200, todos)
	})

	// Get by ID (only if owned by user)
	route.GET("/:id", func(c *gin.Context) {
		var todo models.Todo
		id := c.Param("id")
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(500, gin.H{"error": "User ID not found in context"})
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(500, gin.H{"error": "User ID is not a uint"})
			return
		}
		if result := db.Where("user_id = ?", uid).First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found or not yours"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			return
		}
		c.JSON(200, todo)
	})

	// Update (only if owned by user)
	route.PUT("/:id", func(c *gin.Context) {
		var todo models.Todo
		id := c.Param("id")
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(500, gin.H{"error": "User ID not found in context"})
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(500, gin.H{"error": "User ID is not a uint"})
			return
		}
		if result := db.Where("user_id = ?", uid).First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found or not yours"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			return
		}

		var updatedTodo models.Todo
		if err := c.ShouldBindJSON(&updatedTodo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		todo.Title = updatedTodo.Title
		todo.Description = updatedTodo.Description
		if result := db.Save(&todo); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to update todo: " + result.Error.Error()})
			return
		}
		c.JSON(200, todo)
	})

	// Delete (only if owned by user)
	route.DELETE("/:id", func(c *gin.Context) {
		var todo models.Todo
		id := c.Param("id")
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(500, gin.H{"error": "User ID not found in context"})
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(500, gin.H{"error": "User ID is not a uint"})
			return
		}
		if result := db.Where("user_id = ?", uid).First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found or not yours"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			return
		}
		db.Delete(&todo)
		c.JSON(200, gin.H{"message": "Todo with ID " + id + " deleted"})
	})
}
