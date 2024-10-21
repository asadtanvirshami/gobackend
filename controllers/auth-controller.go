package controllers

import (
	"net/http"
	"os"
	"time"
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"your-app/initializers"
	"your-app/models"
	"google.golang.org/api/idtoken"
	"your-app/utils"
	"github.com/google/uuid"
)

func Signup(c *gin.Context) {
	// Get request body


	var body struct {
		FirstName    string `json:"firstName" binding:"required"`
		LastName    string `json:"lastName" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		fmt.Println(err)

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
		FirstName: body.FirstName,
		LastName:  body.LastName,
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
		"success":true,
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
		"token":   tokenString,
	})
}

// GoogleLogin handles login using Google OAuth 2.0
func GoogleLogin(c *gin.Context) {
	// Request body structure
	var body struct {
		Token string `json:"token" binding:"required"` 
	}

	// Bind JSON body
	if err := c.Bind(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Validate the Google ID token
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	ctx := context.Background()

	payload, err := idtoken.Validate(ctx, body.Token, clientID)
	if err != nil {
		utils.Respond(c, http.StatusUnauthorized, gin.H{
			"error": "Invalid Google token",
		})
		return
	}

	// Extract user information from Google token
	email := payload.Claims["email"].(string)
	fname := payload.Claims["given_name"].(string)
	lname := payload.Claims["family_name"].(string)

	// Check if the user already exists in the database
	var user models.User
	initializers.DB.First(&user, "email = ?", email)

	// If the user doesn't exist, create a new user
	if user.ID == uuid.Nil {
		user = models.User{
			ID: uuid.New(), 
			Email:    email,
			FirstName: fname, 
			LastName: lname, 
			Password: "",  
		}

		// Save the new user in the database
		result := initializers.DB.Create(&user)
		if result.Error != nil {
			utils.Respond(c, http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}
	}

	// Generate a JWT token for the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"id": user.ID,
		"fname": user.FirstName,
		"lname": user.LastName,
		"blocked":user.Blocked,
		"email":user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), 
	})

	// Sign the token
	secretKey := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Set cookie and return token in the response
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true) // 30-day cookie

	// Respond with success
	utils.Respond(c, http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"token":   tokenString,
		"success":true,
	})
}

func Validate(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	utils.Respond(c, http.StatusOK, gin.H{
		"message": "User validated successfully",
		"success":true,
		"user":    user,
	})
}
