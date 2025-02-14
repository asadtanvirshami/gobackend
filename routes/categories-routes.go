package routes

import (
	"your-app/controllers"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categoryRoutes := router.Group("/public/api/v1/category")
	{
		categoryRoutes.GET("/get", controllers.GetCategories)
		categoryRoutes.POST("/create", controllers.CreateCategory)
	}
}
