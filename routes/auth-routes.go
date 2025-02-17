// routes/user_routes.go
package routes

import (
	"your-app/controllers"
	"your-app/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/public/api/v1/user")
	{
		userGroup.POST("/auth/signup", controllers.Signup)
		userGroup.POST("/auth/signin", controllers.Signin)
		userGroup.POST("/auth/otp-verification", controllers.OTPVerification)
		userGroup.POST("/auth/resend-otp", controllers.ResendOTP)
		userGroup.POST("/auth/account-recovery", controllers.AccountRecovery)
		userGroup.POST("/auth/reset-password", controllers.ResetPassword)
		userGroup.POST("/auth/google-signin", controllers.GoogleLogin)
		userGroup.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)
	}
}
