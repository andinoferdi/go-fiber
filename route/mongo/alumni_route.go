package route

import (
	service "go-fiber/app/service/mongo"
	middleware "go-fiber/middleware/mongo"

	"github.com/gofiber/fiber/v2"
)

func AlumniRoutes(app *fiber.App, alumniService *service.AlumniService, authService *service.AuthService) {
	api := app.Group("/go-fiber-mongo")

	api.Post("/login", func(c *fiber.Ctx) error {
		return authService.LoginService(c)
	})

	api.Post("/alumni/check/:key", func(c *fiber.Ctx) error {
		return alumniService.CheckAlumniService(c)
	})

	protected := api.Group("", middleware.AuthRequired())

	protected.Get("/profile", func(c *fiber.Ctx) error {
		return authService.GetProfileService(c)
	})

	alumni := protected.Group("/alumni")
	alumni.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return alumniService.GetAllAlumniService(c)
	})
	alumni.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return alumniService.GetAlumniByIDService(c)
	})
	alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return alumniService.CreateAlumniService(c)
	})
	alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return alumniService.UpdateAlumniService(c)
	})
	alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return alumniService.DeleteAlumniService(c)
	})
}
