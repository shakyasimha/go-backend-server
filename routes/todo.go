package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

func connectDB() *gorm.DB {
	// Connect to the SQLite DB
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database.")
	}

	// Automigrate the Todo model to create the table
	db.AutoMigrate(&Todo{})

	// Returns db
	return db
}

func TodoRoutes(router *gin.Engine, db *gorm.DB) {
	// Routes go here
	route := router.Group("/todos")

	// APIs go here
	// Route to create a new Todo
	route.POST("/todos", func(c *gin.Context) {
		var todo Todo

		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		// Save the Todo to db
		if result := db.Create(&todo); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create todo: " + result.Error.Error()})
			return
		}

		c.JSON(200, todo)
	})

	// Route to get all Todo
	route.GET("/todos", func(c *gin.Context) {
		var todos []Todo

		// Retrieve all Todos from database
		if result := db.Find(&todos); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve todos: " + result.Error.Error()})
			return
		}

		c.JSON(200, todos)
	})

	// Router to get Todo by id
	route.GET("/todos/:id", func(c *gin.Context) {
		var todo Todo

		// Getting id
		id := c.Param("id")

		// Retrieve data by id
		if result := db.First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found."})
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			return
		}

		c.JSON(200, todo)
	})

	// Route to edit todo
	router.PUT("/todos", func(c *gin.Context) {
		var todo Todo
		id := c.Param("id")

		if result := db.First(&todo, id); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "Todo not found."})
			}
			c.JSON(500, gin.H{"error": "Failed to retrieve todo: " + result.Error.Error()})
			return
		}

		var updatedTodo Todo
		if err := c.ShouldBindJSON(&updatedTodo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		// Update the todo in database
		todo.Title = updatedTodo.Title
		todo.Description = updatedTodo.Description

		if result := db.Save(&todo); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to update todo: " + result.Error.Error()})
			return
		}

		c.JSON(200, todo)
	})

	// Route to delete a Todo by ID
	router.DELETE("/todos/:id", func(c *gin.Context) {
		var todo Todo
		id := c.Param("id")

		// Retrieve todo from DB
		if result := db.First(&todo, id); result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found."})
			return
		}

		db.Delete(&todo)

		c.JSON(200, gin.H{"message": fmt.Sprintf("Todo with ID %s is deleted.", id)})
	})
}
