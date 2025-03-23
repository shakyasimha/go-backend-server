package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/middleware"
)

type UserController struct{}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	userID := c.Param("id")

	c.JSON(200, gin.H{
		"id":    userID,
		"name":  "John Doe",
		"email": "john@example.com",
	})
}

// SetupRoutes configures all routes
func SetupRoutes(router *gin.Engine) {
	public := router.Group("/public")

	userController := &UserController{}

	public.Use(middleware.LoggerMiddleware())
	{
		public.GET("/info", func(c *gin.Context) {
			c.String(200, "Public information")
		})

		public.GET("/products", func(c *gin.Context) {
			c.String(200, "Public product list")
		})
	}

	private := router.Group("/private")
	private.Use(middleware.LoggerMiddleware())
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/data", func(c *gin.Context) {
			c.String(200, "Private data accessible after authentication.")
		})

		private.POST("/create", func(c *gin.Context) {
			c.String(200, "Create a new resource")
		})
	}

	userRouter := router.Group("/users")
	private.Use(middleware.LoggerMiddleware())
	{
		userRouter.GET("/users/:id", userController.GetUserInfo)
	}

}
