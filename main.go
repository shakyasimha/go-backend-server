package main

import (
	"github.com/shakyasimha/go-backend-server/routes"
	"github.com/shakyasimha/go-backend-server/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	// Database connection
	db := utils.ConnectDB("todos.db")

	// Set up routes
	routes.SetupRoutes(router)
	routes.TodoRoutes(router, db)

	// Run server
	router.Run(":8080")
}
