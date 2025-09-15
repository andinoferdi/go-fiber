package config

import (
	"database/sql"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Alumni Management API",
	})
	app.Use(middleware.LoggerMiddleware)
	return app
}
