package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/routes"
)

func main() {
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":8080")
}
