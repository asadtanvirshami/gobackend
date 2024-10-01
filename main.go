package main 

//imports
import (
	"fmt"
	"os"
	"your-app/initializers"
	"your-app/routes"
	"github.com/gin-gonic/gin"
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

	// Setup routes
	routes.UserRoutes(router)

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
