package controllers

import (
	"net/http"
	"your-app/initializers"
	"your-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCommunities(c *gin.Context) {
	var communities []models.Community
	initializers.DB.Find(&communities)
	c.JSON(http.StatusOK, communities)
}

func CreateCommunity(c *gin.Context) {
	var community models.Community
	if err := c.ShouldBindJSON(&community); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	community.ID = uuid.New()
	initializers.DB.Create(&community)
	c.JSON(http.StatusCreated, community)
}

func GetCommunityByID(c *gin.Context) {
	id := c.Param("id")
	var community models.Community
	if err := initializers.DB.First(&community, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Community not found"})
		return
	}
	c.JSON(http.StatusOK, community)
}
