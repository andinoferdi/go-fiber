package route

import (
	"database/sql"
	postgre "go-fiber/app/service/postgre"
	middleware "go-fiber/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func AlumniRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/go-fiber-postgre")

	api.Post("/login", func(c *fiber.Ctx) error {
		return postgre.LoginService(c, db)
	})

	protected := api.Group("", middleware.AuthRequired())
	
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return postgre.GetProfileService(c, db)
	})

	alumni := protected.Group("/alumni")
	alumni.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.GetAllAlumniService(c, db)
	})
	alumni.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.GetAlumniByIDService(c, db)
	})
	alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.CreateAlumniService(c, db)
	})
	alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.UpdateAlumniService(c, db)
	})
	alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.DeleteAlumniService(c, db)
	})
	alumni.Post("/check/:key", func(c *fiber.Ctx) error {
		return postgre.CheckAlumniService(c, db)
	})

	roles := protected.Group("/roles")
	roles.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.GetAllRolesService(c, db)
	})
	roles.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return postgre.GetRoleByIDService(c, db)
	})
	roles.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.CreateRoleService(c, db)
	})
	roles.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.UpdateRoleService(c, db)
	})
	roles.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return postgre.DeleteRoleService(c, db)
	})
}
