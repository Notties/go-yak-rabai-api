// config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"context"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"yak.rabai/models"
)

// Global variables for DB and Redis clients
var DB *gorm.DB
var RDB *redis.Client

func init() {
	// Load environment variables from .env file (if it exists)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// ConnectDatabase connects to PostgreSQL using environment variables
func ConnectDatabase() {
	// Get environment variables for PostgreSQL connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construct the PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable query logging (optional)
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db
	fmt.Println("Connected to PostgreSQL database!")
	MigrateDB()
	fmt.Println("Migrations to PostgreSQL database!")
}

// MigrateDB runs the migrations to create tables for the models
func MigrateDB() {
	// Automatically migrate the schema
	if err := DB.AutoMigrate(
		&models.User{},
		&models.ChatRoom{},
		&models.Rating{},
		// Add other models here as needed
	); err != nil {
		log.Fatal("Migration failed: ", err)
	}
	fmt.Println("Migration completed successfully")
}

// ConnectRedis connects to Redis using environment variables
func ConnectRedis() {
	// Get Redis host and port from environment variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	// Check Redis connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	RDB = client
	fmt.Println("Connected to Redis!")
}
