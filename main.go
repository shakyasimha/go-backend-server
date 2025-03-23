package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/routes"
	"github.com/shakyasimha/go-backend-server/todo"
)

func main() {
	router := gin.Default()

	// Root directory here
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	// Routes for routes setup here
	routes.SetupRoutes(router)

	// Routes for todo setup here
	todo.SetupRoutes(router)

	// Server port run here
	router.Run(":8080")
}
