package main

import (
	"github.com/shakyasimha/go-backend-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	// Set up routes
	routes.SetupRoutes(router)
	routes.TodoRoutes(router)

	// Run server
	router.Run(":8080")
}
