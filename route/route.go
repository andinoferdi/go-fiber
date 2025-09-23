package route

import (
	"database/sql"
	"go-fiber/app/service"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/go-fiber")

	api.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	protected := api.Group("", middleware.AuthRequired())
	
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	})

	rolesRoutes := protected.Group("/roles")
	rolesRoutes.Get("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetAllRolesService(c, db)
	})
	rolesRoutes.Get("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetRoleByIDService(c, db)
	})

	alumniRoutes := protected.Group("/alumni")
	alumniRoutes.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})
	alumniRoutes.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})
	alumniRoutes.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})
	alumniRoutes.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})
	alumniRoutes.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})

	pekerjaanRoutes := protected.Group("/pekerjaan")
	pekerjaanRoutes.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanAlumniService(c, db)
	})
	pekerjaanRoutes.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByIDService(c, db)
	})
	pekerjaanRoutes.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanByAlumniIDService(c, db)
	})
	pekerjaanRoutes.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanAlumniService(c, db)
	})
	pekerjaanRoutes.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniService(c, db)
	})
	pekerjaanRoutes.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeletePekerjaanAlumniService(c, db)
	})
	pekerjaanRoutes.Put("/soft-delete/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.SoftDeletePekerjaanAlumniService(c, db)
	})
}
