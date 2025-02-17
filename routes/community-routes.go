package routes

import (
	"your-app/controllers"

	"github.com/gin-gonic/gin"
)

func CommunityRoutes(r *gin.Engine) {
	community := r.Group("/communities")
	{
		community.GET("/", controllers.GetCommunities)
		community.POST("/", controllers.CreateCommunity)
		community.GET("/:id", controllers.GetCommunityByID)
	}
}
