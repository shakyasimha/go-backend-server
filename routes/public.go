package routes

import (
	"go-backend-server/utils" // Updated to match module

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(200, gin.H{"id": userID, "name": "John Doe", "email": "john@example.com"})
}

func SetupRoutes(router *gin.Engine) {
	public := router.Group("/public")
	userController := &UserController{}
	public.Use(utils.LoggerMiddleware())
	{
		public.GET("/info", func(c *gin.Context) { c.String(200, "Public information") })
		public.GET("/products", func(c *gin.Context) { c.String(200, "Public product list") })
	}

	private := router.Group("/private")
	private.Use(utils.LoggerMiddleware())
	private.Use(utils.AuthMiddleware())
	{
		private.GET("/data", func(c *gin.Context) { c.String(200, "Private data...") })
		private.POST("/create", func(c *gin.Context) { c.String(200, "Create a new resource") })
	}

	userRouter := router.Group("/users")
	userRouter.Use(utils.LoggerMiddleware())
	{
		userRouter.GET("/:id", userController.GetUserInfo)
	}
}
