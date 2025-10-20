package route

import (
	service "go-fiber/app/service/mongo"
	middleware "go-fiber/middleware/mongo"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app *fiber.App, pekerjaanService *service.PekerjaanAlumniService) {
	api := app.Group("/go-fiber-mongo")
	
	pekerjaan := api.Group("/pekerjaan", middleware.AuthRequired())
	pekerjaan.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return pekerjaanService.GetAllPekerjaanAlumniService(c)
	})
	pekerjaan.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return pekerjaanService.GetPekerjaanAlumniByIDService(c)
	})
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return pekerjaanService.GetPekerjaanAlumniByAlumniIDService(c)
	})
	pekerjaan.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return pekerjaanService.CreatePekerjaanAlumniService(c)
	})
	pekerjaan.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return pekerjaanService.UpdatePekerjaanAlumniService(c)
	})
	pekerjaan.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return pekerjaanService.DeletePekerjaanAlumniService(c)
	})
}
