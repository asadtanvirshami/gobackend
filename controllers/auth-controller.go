package controllers

import (
	"net/http"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"your-app/initializers"
	"your-app/models"
	"your-app/utils"
	"github.com/google/uuid"
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

func Signin(c *gin.Context) {
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

	// Find the user by email
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	// Check if the user exists by comparing the UUID to uuid.Nil
	if user.ID == uuid.Nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	
	// Compare and send in pass with saved hash pass
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	// Generate Jwt Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign the token with the secret key
	secretKey := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey)) // Convert secret key to []byte
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Respond with success
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization",tokenString, 3600 * 24 * 30, "", "", false, true)	

	c.JSON(http.StatusOK, gin.H{
		"success":true,
		"message": "User logged in successfully",
	})
}

func Validate(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	utils.Respond(c, http.StatusOK, gin.H{
		"message": "User validated successfully",
		"user":    user,
	})
}
