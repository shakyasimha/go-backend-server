package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/routes"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	routes.SetupRoutes(router)
	router.Run(":8080")
}
