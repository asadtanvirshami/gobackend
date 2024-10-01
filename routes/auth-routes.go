// routes/user_routes.go
package routes

import (
	"your-app/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/public/api/v1/user")
	{
		userGroup.POST("/register", controllers.Signup)
	}
}
