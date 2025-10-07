package main

import (
	"go-fiber/config"
	"go-fiber/database"
	"go-fiber/middleware"
	"go-fiber/route"
	"log"
	"os"
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()
	defer db.Close()
	
	app := config.NewApp(db)
	app.Use(middleware.LoggerMiddleware)
	
	route.AlumniRoutes(app, db)
	route.PekerjaanRoutes(app, db)
	
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
