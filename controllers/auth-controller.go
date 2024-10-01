package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"your-app/initializers"
	"your-app/models"
	"your-app/utils"
)

func Signup(c *gin.Context) {
	// Get request body
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Respond with success
	utils.Respond(c, http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}
