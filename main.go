package main 

//imports
import (
	"fmt"
	"os"
	"your-app/initializers"
	"your-app/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"log"
)

//init func for env initialzation
func init() {
	initializers.LoadEnv()
	initializers.DBConnection()
	initializers.SyncDB()
}

//main func for server connection
func main() {
	// Initialize Gin router
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust this to your front-end domain in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}

	// Use CORS middleware
	router.Use(cors.New(corsConfig))


	// Setup routes
	routes.UserRoutes(router)

	// *****Starting Server*****
	fmt.Println("||==Starting Server==||")
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
