package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

func TodoRoutes(router *gin.Engine, db *gorm.DB) {
	route := router.Group("/todos")

	// Create
	route.POST("/", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}
		if result := db.Create(&todo); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create todo: " + result.Error.Error()})
			return
		}
		c.JSON(200, todo)
	})

	// List
	route.GET("/", func(c *gin.Context) {
		var todos []Todo
		if result := db.Find(&todos); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve todos: " + result.Error.Error()})
			return
		}
		c.JSON(200, todos)
	})

	// Get by ID
	route.GET("/:id", func(c *gin.Context) {
		var todo Todo
		id := c.Param("id")
		if result := db.First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			return
		}
		c.JSON(200, todo)
	})

	// Update
	route.PUT("/:id", func(c *gin.Context) {
		var todo Todo
		id := c.Param("id")
		if result := db.First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			return
		}
		var updatedTodo Todo
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

	// Delete
	route.DELETE("/:id", func(c *gin.Context) {
		var todo Todo
		id := c.Param("id")
		if result := db.First(&todo, id); result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}
		db.Delete(&todo)
		c.JSON(200, gin.H{"message": "Todo with ID " + id + " is deleted"})
	})
}
