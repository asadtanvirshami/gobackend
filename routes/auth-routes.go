// routes/user_routes.go
package routes

import (
	"your-app/controllers"
	"your-app/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/public/api/v1/user")
	{
		userGroup.POST("/register", controllers.Signup)
		userGroup.POST("/signin", controllers.Signin)
		userGroup.GET("/validate",middleware.RequireAuth, controllers.Validate)
	}
}
