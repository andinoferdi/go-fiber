package main

import (
	"go-fiber/config"
	"go-fiber/database"
	"go-fiber/middleware"
	routepostgre "go-fiber/route/postgre"
	"log"
	"os"
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()
	defer db.Close()
	
	app := config.NewApp(db)
	app.Use(middleware.LoggerMiddleware)
	
	routepostgre.AlumniRoutes(app, db)
	routepostgre.PekerjaanRoutes(app, db)
	
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
