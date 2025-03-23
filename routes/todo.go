package routes

import (
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
	route := router.Group("/todo")

	// APIs go here
	// Route to create a new Todo
	route.POST("/todo", func(c *gin.Context) {
		var todo Todo

		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		// Save the Todo to db
		if result := db.Create(&todo); result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create todo." + result.Error.Error()})
			return 
		}

		c.JSON(200, todo)
	})

	// Route to get all Todo 
	route.GET("/todo", func(c *gin.Context) {
		var todos []Todo 
		
		// Retrieve all Todos from database 
		if result := db.Find(&todos); result.Error != nil {
			c.JSON(500, gin.H{
				"error": "Failed to retrieve todos: " + result.Error.Error()
			})
		}

		c.JSON(200, todos)
	})
}
