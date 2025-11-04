package main

import (
	repositorymongo "go-fiber/app/repository/mongo"
	servicemongo "go-fiber/app/service/mongo"
	"go-fiber/config"
	configmongo "go-fiber/config/mongo"

	"go-fiber/database"
	_ "go-fiber/docs"
	"go-fiber/middleware"
	routemongo "go-fiber/route/mongo"

	routepostgre "go-fiber/route/postgre"
	"log"
	"os"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Alumni Management System API - MongoDB
// @version 1.0
// @description API untuk mengelola data alumni, pekerjaan, dan file menggunakan MongoDB dengan JWT Authentication
// @host localhost:3000
// @BasePath /go-fiber-mongo
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format: Bearer {token}

// @tag.name 1. Authentication
// @tag.description Operasi autentikasi dan profil pengguna

// @tag.name 2. Alumni
// @tag.description Operasi CRUD untuk data alumni

// @tag.name 3. Pekerjaan Alumni
// @tag.description Operasi CRUD untuk data pekerjaan alumni

// @tag.name 4. Files
// @tag.description Operasi upload dan manajemen file

func main() {
	config.LoadEnv()
	
	postgresDB := database.ConnectDB()
	defer postgresDB.Close()
	
	mongoDB := database.ConnectMongoDB()
	
	// Run MongoDB migrations
	if err := database.RunMigrations(mongoDB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	
	
	alumniRepo := repositorymongo.NewAlumniRepository(mongoDB)
	alumniService := servicemongo.NewAlumniService(alumniRepo)
	
	authService := servicemongo.NewAuthService(alumniRepo)
	
	pekerjaanRepo := repositorymongo.NewPekerjaanAlumniRepository(mongoDB)
	pekerjaanService := servicemongo.NewPekerjaanAlumniService(pekerjaanRepo)
	
	fileRepo := repositorymongo.NewFileRepository(mongoDB)
	fileService := servicemongo.NewFileService(fileRepo, alumniRepo, "./uploads")
	
	app := configmongo.NewApp()
	app.Use(middleware.LoggerMiddleware)
	
	routepostgre.AlumniRoutes(app, postgresDB)
	routepostgre.PekerjaanRoutes(app, postgresDB)
	
	routemongo.AlumniRoutes(app, alumniService, authService)
	routemongo.PekerjaanRoutes(app, pekerjaanService)
	routemongo.FileRoutes(app, fileService)
	
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
