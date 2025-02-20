package initializers

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB and Redis client instances
var (
	DB    *gorm.DB
	Redis *redis.Client
	Ctx   = context.Background() // Context for Redis operations
)

// DBConnection initializes PostgreSQL and Redis
func DBConnection() {
	// Load environment variables
	connectionString := os.Getenv("CONNECTION_STRING")
	redisAddr := os.Getenv("REDIS_ADDR") // Example: "localhost:6379"
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB") // Example: "0"

	// Connect to PostgreSQL
	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}
	log.Println("✅ PostgreSQL connected successfully")

	// Connect to Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // No password if empty
		DB:       0,             // Default DB
	})

	// Test Redis connection
	_, err = Redis.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("❌ Failed to connect to Redis:", err)
	}
	log.Println("✅ Redis connected successfully")
}
