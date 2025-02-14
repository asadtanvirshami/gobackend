package controllers

import (
	"net/http"
	// "os"
	// "time"
	// "context"
	// "fmt"

	"your-app/initializers"
	"your-app/models"
	"your-app/utils"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

func CreateCategory(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Create the user
	category := models.Category{
		Name: body.Name,
	}

	result := initializers.DB.Create(&category)
	if result.Error != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "hit is successfully",
		"data":    body,
	})
}
func GetCategories(c *gin.Context) {
	var categories []models.Category

	result := initializers.DB.Find(&categories)
	if result.Error != nil {
		utils.Respond(c, http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Categories fetched successfully",
		"data":    categories,
	})
}
