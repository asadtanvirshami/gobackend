package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"your-app/initializers"
	"your-app/models"
	"github.com/google/uuid"
)

func RequireAuth(c *gin.Context) {
	// Get the cookie from the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, c.AbortWithError(http.StatusUnauthorized, err)
		}

		// Return the secret key for validating the token
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract claims and check validity
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the user based on the subject (sub) claim
		var user models.User
		initializers.DB.First(&user, "id = ?", claims["sub"])

		if user.ID == uuid.Nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach the user to the context
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
