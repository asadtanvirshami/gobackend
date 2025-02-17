package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"your-app/initializers"
	"your-app/models"
	"your-app/services"
	"your-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

func generateOTP() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return fmt.Sprintf("%06d", r.Intn(1000000)) // 6-digit OTP
}

func Signup(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Get file from request
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{"error": "Failed to get image"})
		return
	}
	defer file.Close()

	uploadDir := "uploads"
	// Ensure the "uploads" directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	timestamp := time.Now().UnixMilli()
	fileName := fmt.Sprintf("%d-%s", timestamp, header.Filename)
	imagePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	_, err = io.Copy(out, file)
	if err != nil {
		out.Close() // Ensure we close the file before responding
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write image file"})
		return
	}

	out.Sync()  // Flush any remaining writes
	out.Close() // Close immediately to release file lock

	// Open the file for reading (for ImageKit upload)
	fileData, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file for reading:", err)
		return
	}

	// Read into buffer
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, fileData)
	if err != nil {
		fileData.Close()
		fmt.Println("Error copying file to buffer:", err)
		return
	}

	fileData.Close() // Close file after reading to release OS lock

	// Upload image
	url, fileId, err := services.UploadImage(imagePath, fileName, "documents")
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate OTP
	otp := generateOTP()
	expiresAt := time.Now().Add(1 * time.Minute).Unix()

	// Create User
	user := models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Password:     string(hashedPassword),
		Image:        &url,
		ImageID:      &fileId,
		OTP:          &otp,
		OTPExpiresAt: &expiresAt,
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	// Send OTP Email
	go services.SendOTPEmail(email, otp)

	// Success Response
	utils.Respond(c, http.StatusOK, gin.H{"success": true, "message": "User created successfully"})

	// Wait for OS buffer clearance (last resort)
	time.Sleep(100 * time.Millisecond)

	// Now we can safely delete the file
	if err := os.Remove(imagePath); err != nil {
		fmt.Println("Error removing file:", err)
	} else {
		fmt.Println("File successfully removed:", imagePath)
	}
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
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User logged in successfully",
		"token":   tokenString,
	})
}

func OTPVerification(c *gin.Context) {
	var body struct {
		OTP string `json:"otp" binding:"required"`
	}

	// Bind the JSON request body
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Retrieve the user with the given OTP from the database
	var user models.User
	if result := initializers.DB.Where("otp = ?", body.OTP).First(&user); result.Error != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Invalid OTP or user not found",
		})
		return
	}

	// Check if OTP is expired
	currentTime := time.Now().Unix()
	if user.OTPExpiresAt != nil && *user.OTPExpiresAt < currentTime {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "OTP has expired",
		})
		return
	}

	// Verify OTP
	if user.OTP != nil && *user.OTP == body.OTP {
		utils.Respond(c, http.StatusOK, gin.H{
			"success": true,
			"message": "OTP verified successfully",
		})
	} else {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Invalid OTP",
		})
	}
}

func ResendOTP(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Bind the JSON request body
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Retrieve the user by email
	var user models.User
	if result := initializers.DB.Where("email = ?", body.Email).First(&user); result.Error != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	// Generate a new OTP and set expiration time (1 minute from now)
	newOTP := generateOTP()                                // Implement the OTP generation logic here
	otpExpiresAt := time.Now().Add(1 * time.Minute).Unix() // Expiry time set to 1 minute from now

	// Update the user's OTP and expiration time
	user.OTP = &newOTP
	user.OTPExpiresAt = &otpExpiresAt
	if err := initializers.DB.Save(&user).Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, gin.H{
			"error": "Failed to update OTP",
		})
		return
	}

	// Send the OTP to the user's email
	go services.SendOTPEmail(body.Email, newOTP)

	utils.Respond(c, http.StatusOK, gin.H{
		"success": true,
		"message": "OTP sent successfully",
	})
}

func AccountRecovery(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Bind the JSON request body
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Find the user by email
	var user models.User
	if result := initializers.DB.Where("email = ?", body.Email).First(&user); result.Error != nil {
		utils.Respond(c, http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	// Generate a new OTP and set expiration time (1 minute from now)
	newOTP := generateOTP()
	otpExpiresAt := time.Now().Add(1 * time.Minute).Unix()
	// Update the user's OTP and expiration time
	user.OTP = &newOTP
	user.OTPExpiresAt = &otpExpiresAt
	if err := initializers.DB.Save(&user).Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, gin.H{
			"error": "Failed to update OTP",
		})
		return
	}

	go services.SendOTPEmail(body.Email, newOTP)

	utils.Respond(c, http.StatusOK, gin.H{
		"success": true,
		"message": "OTP sent successfully",
	})
}

func ResetPassword(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var user models.User
	if err := initializers.DB.Save(&user).Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, gin.H{
			"error": "Failed to update OTP",
		})
		return
	}

	user.Password = body.Password
	if err := initializers.DB.Save(&user).Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, gin.H{
			"error": "Failed to update password",
		})
		return
	}

	utils.Respond(c, http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset successfully",
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
			ID:        uuid.New(),
			Email:     email,
			FirstName: fname,
			LastName:  lname,
			Password:  "",
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
		"sub":     user.ID,
		"id":      user.ID,
		"fname":   user.FirstName,
		"lname":   user.LastName,
		"blocked": user.Blocked,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
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
		"success": true,
	})
}

func Validate(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	utils.Respond(c, http.StatusOK, gin.H{
		"message": "User validated successfully",
		"success": true,
		"user":    user,
	})
}
