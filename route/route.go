package route

import (
	"database/sql"
	"go-fiber/app/service"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	// API Group
	api := app.Group("/go-fiber")

	// Public Routes - Login
	api.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	// Protected Routes - Require Authentication
	protected := api.Group("", middleware.AuthRequired())
	
	// Profile endpoint - accessible by all authenticated users
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	})

	// Alumni Routes with Access Control
	alumniRoutes := protected.Group("/alumni")
	// GET endpoints - accessible by both admin and user
	alumniRoutes.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})
	alumniRoutes.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})
	// CUD endpoints - admin only
	alumniRoutes.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})
	alumniRoutes.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})
	alumniRoutes.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})

	// Pekerjaan Alumni Routes with Access Control
	pekerjaanRoutes := protected.Group("/pekerjaan")
	// GET all pekerjaan - accessible by both admin and user
	pekerjaanRoutes.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanAlumniService(c, db)
	})
	// GET pekerjaan by ID - accessible by both admin and user
	pekerjaanRoutes.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByIDService(c, db)
	})
	// GET pekerjaan by alumni ID - admin only per requirements
	pekerjaanRoutes.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanByAlumniIDService(c, db)
	})
	// CUD endpoints - admin only
	pekerjaanRoutes.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanAlumniService(c, db)
	})
	pekerjaanRoutes.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniService(c, db)
	})
	pekerjaanRoutes.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeletePekerjaanAlumniService(c, db)
	})
}
