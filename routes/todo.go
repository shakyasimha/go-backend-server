// routes/todo.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/models"
	"github.com/shakyasimha/go-backend-server/utils"
	"gorm.io/gorm"
)

func TodoRoutes(router *gin.Engine) {
	todo := models.NewTodo()
	db := todo.ConnectDB()

	route := router.Group("/todos", utils.AuthMiddleware(db))

	route.POST("/", func(c *gin.Context) {
		var todo models.Todo

		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		userID, _ := c.Get("user_id")
		todo.UserID = userID.(uint)

		if result := db.Create(&todo); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create todo: " + result.Error.Error()})
			return
		}
		c.JSON(200, todo)
	})

	route.GET("/", func(c *gin.Context) {
		var todos []models.Todo

		userID, _ := c.Get("user_id")

		if result := db.Where("user_id = ?", userID).Find(&todos); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve todos: " + result.Error.Error()})
			return
		}
		c.JSON(200, todos)
	})

	route.GET("/:id", func(c *gin.Context) {
		var todo models.Todo

		id := c.Param("id")
		userID, _ := c.Get("user_id")

		if result := db.Where("user_id = ?", userID).First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found or not yours"})
			} else {
				c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			}
			return
		}
		c.JSON(200, todo)
	})

	route.PUT("/:id", func(c *gin.Context) {
		var todo models.Todo

		id := c.Param("id")
		userID, _ := c.Get("user_id")

		if result := db.Where("user_id = ?", userID).First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found or not yours"})
			} else {
				c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			}
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

	route.DELETE("/:id", func(c *gin.Context) {
		var todo models.Todo

		id := c.Param("id")
		userID, _ := c.Get("user_id")

		if result := db.Where("user_id = ?", userID).First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found or not yours"})
			} else {
				c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			}
			return
		}

		db.Delete(&todo)
		c.JSON(200, gin.H{"message": "Todo with ID " + id + " deleted"})
	})
}
