package route

import (
	"database/sql"
	"go-fiber/app/service"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/go-fiber")
	protected := api.Group("", middleware.AuthRequired())

	pekerjaan := protected.Group("/pekerjaan")
	
	pekerjaan.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanAlumniService(c, db)
	})
	
	// Tugas UTS 1 Get All Soft Delet
	pekerjaan.Get("/soft-delete", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllSoftDeletedPekerjaanAlumniService(c, db)
	})
	
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByAlumniIDService(c, db)
	})
	pekerjaan.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByIDService(c, db)
	})
	
	pekerjaan.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanAlumniService(c, db)
	})
	pekerjaan.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniService(c, db)
	})

	// Tugas sebelumnya soft-delete
	pekerjaan.Put("/soft-delete/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.SoftDeletePekerjaanAlumniService(c, db)
	})

	// Tugas UTS 2 Soft Delete Restore
	pekerjaan.Put("/soft-delete-restore/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.RestorePekerjaanAlumniService(c, db)
	})

	// Tugas UTS 3 Hard Delete
	pekerjaan.Delete("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return service.HardDeletePekerjaanAlumniService(c, db)
	})
}
