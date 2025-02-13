package controllers

import (
		"net/http"
		// "os"
		// "time"
		// "context"
		// "fmt"

		"github.com/gin-gonic/gin"
		"your-app/utils"
		// "your-app/initializers"
		// "your-app/models"
		// "github.com/google/uuid"
	)

	func CreateCategory (c *gin.Context){
		var body struct {
			Name    string `json:"name" binding:"required"`
		}
	
		if err := c.Bind(&body); err != nil {
			utils.Respond(c, http.StatusBadRequest, gin.H{
				"error": "Failed to read body",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success":true,
			"message": "hit is successfully",
			"data": body,
		})
	}