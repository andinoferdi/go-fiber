package route

import (
	"database/sql"
	postgre "go-fiber/app/service/postgre"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/go-fiber")
	protected := api.Group("", middleware.AuthRequired())

	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.GetAllPekerjaanAlumniService(c, db)
	})
	pekerjaan.Get("/trash", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.GetSoftDeletedPekerjaanAlumniService(c, db)
	})
	pekerjaan.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.GetPekerjaanAlumniByIDService(c, db)
	})
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.GetPekerjaanAlumniByAlumniIDService(c, db)
	})
	pekerjaan.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.CreatePekerjaanAlumniService(c, db)
	})
	pekerjaan.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.UpdatePekerjaanAlumniService(c, db)
	})
	pekerjaan.Put("/soft-delete/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.SoftDeletePekerjaanAlumniService(c, db)
	})
	pekerjaan.Put("/soft-delete-restore/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.RestorePekerjaanAlumniService(c, db)
	})
	pekerjaan.Delete("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.HardDeletePekerjaanAlumniService(c, db)
	})
}