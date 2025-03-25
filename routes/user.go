<<<<<<< HEAD
// struct User
// name, email, password,

// struct Token
//
//

// CRUD

// jwt token 
=======
package routes
>>>>>>> 2fe90be3a3a16687520c5364e5b7f7f61ab4b110
import (
	"github.com/gin-gonic/gin"
	"github.com/shakyasimha/go-backend-server/models"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.Engine) {
	user := models.User{}
	db := user.ConnectDB() // Connects the db for orm task s

	// Declaring the routes here 
	route := router.Group("/users")

	/*
		CRUD requests go here
	*/	
	// For creating a new user
	route.POST("/create", func(c *gin.Context) {
		var user models.User

		// For checking if the data is valid JSON or not
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		// Returns the error if data is not created
		if results := db.Create(&user); results.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create user: " + results.Error.Error()})
			return
		}

		// Should create an authentication API that returns token or does something to authenticate the user	
		// Main task here: how to authenticate the user and how to allow them access to todo lists
		c.JSON(200, gin.H{"user": user})
	})

	//	For logging in the user 
	route.
}