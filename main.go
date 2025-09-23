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
	config.InitLogger()
	logger := config.GetLogger()
	middleware.InitRequestLogger()
	config.LoadEnv()
	logger.Println("Environment variables loaded")
	db := database.ConnectDB()
	logger.Println("Database connected successfully")
	app := config.NewApp(db)
	logger.Println("Fiber app initialized")
	route.RegisterRoutes(app, db)
	logger.Println("Routes registered successfully")
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	logger.Printf("Starting server on port %s", port)
	fmt.Printf("🚀 Alumni API is running on http://localhost:%s\n", port)
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
