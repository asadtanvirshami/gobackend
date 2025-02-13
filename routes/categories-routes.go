package routes

import (
	"your-app/controllers"
	"your-app/middleware"
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categoryRoutes := router.Group("/public/api/v1/category")
	{
		categoryRoutes.GET("/get",middleware.RequireAuth, controllers.Validate)
		categoryRoutes.POST("/create", controllers.CreateCategory)
		categoryRoutes.PUT("/update", controllers.Signin)
		categoryRoutes.DELETE("/delete", controllers.GoogleLogin)
	}
}
