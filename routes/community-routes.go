package routes

import (
	"your-app/controllers"

	"github.com/gin-gonic/gin"
)

func CommunityRoutes(r *gin.Engine) {
	community := r.Group("/public/api/v1/community")
	{
		community.GET("/get", controllers.GetCommunities)
		community.POST("/create", controllers.CreateCommunity)
		community.GET("/:id", controllers.GetCommunityByID)
	}
}
