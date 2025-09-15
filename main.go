package main

import (
	"fmt"
	"go-fiber/config"
	"go-fiber/database"
	"go-fiber/middleware"
	"go-fiber/route"
	"os"
)

func main() {
	// Initialize logger first
	config.InitLogger()
	logger := config.GetLogger()
	
	// Initialize request logger middleware
	middleware.InitRequestLogger()
	
	// Load environment variables
	config.LoadEnv()
	logger.Println("Environment variables loaded")
	
	// Connect to database
	db := database.ConnectDB()
	logger.Println("Database connected successfully")
	
	// Initialize Fiber app
	app := config.NewApp(db)
	logger.Println("Fiber app initialized")
	
	// Register routes
	route.RegisterRoutes(app, db)
	logger.Println("Routes registered successfully")
	
	// Get port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	
	// Log server startup
	logger.Printf("Starting server on port %s", port)
	fmt.Printf("🚀 Alumni API is running on http://localhost:%s\n", port)
	
	// Start server
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
