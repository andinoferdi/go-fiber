package config

import (
	"github.com/gofiber/fiber/v2"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024, // 2MB limit for file uploads
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Internal server error: " + err.Error(),
			})
		},
	})
	
	// Serve static files from uploads directory
	app.Static("/uploads", "./uploads")
	
	return app
}

