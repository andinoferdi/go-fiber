package main

import (
	repositorymongo "go-fiber/app/repository/mongo"
	servicemongo "go-fiber/app/service/mongo"
	"go-fiber/config"
	configmongo "go-fiber/config/mongo"

	"go-fiber/database"
	"go-fiber/middleware"
	routemongo "go-fiber/route/mongo"

	routepostgre "go-fiber/route/postgre"
	"log"
	"os"
)

func main() {
	config.LoadEnv()
	
	postgresDB := database.ConnectDB()
	defer postgresDB.Close()
	
	mongoDB := database.ConnectMongoDB()
	
	// Run MongoDB migrations
	if err := database.RunMigrations(mongoDB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	
	roleRepo := repositorymongo.NewRoleRepository(mongoDB)
	roleService := servicemongo.NewRoleService(roleRepo)
	
	alumniRepo := repositorymongo.NewAlumniRepository(mongoDB)
	alumniService := servicemongo.NewAlumniService(alumniRepo)
	
	authService := servicemongo.NewAuthService(alumniRepo)
	
	pekerjaanRepo := repositorymongo.NewPekerjaanAlumniRepository(mongoDB)
	pekerjaanService := servicemongo.NewPekerjaanAlumniService(pekerjaanRepo)
	
	fileRepo := repositorymongo.NewFileRepository(mongoDB)
	fileService := servicemongo.NewFileService(fileRepo, "./uploads")
	
	app := configmongo.NewApp()
	app.Use(middleware.LoggerMiddleware)
	
	routepostgre.AlumniRoutes(app, postgresDB)
	routepostgre.PekerjaanRoutes(app, postgresDB)
	
	routemongo.AlumniRoutes(app, alumniService, authService)
	routemongo.RoleRoutes(app, roleService)
	routemongo.PekerjaanRoutes(app, pekerjaanService)
	routemongo.FileRoutes(app, fileService)
	
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
