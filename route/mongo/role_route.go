package route

import (
	service "go-fiber/app/service/mongo"
	middleware "go-fiber/middleware/mongo"

	"github.com/gofiber/fiber/v2"
)

func RoleRoutes(app *fiber.App, roleService *service.RoleService) {
	api := app.Group("/go-fiber-mongo")
	
	roles := api.Group("/roles", middleware.AuthRequired())
	roles.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return roleService.GetAllRolesService(c)
	})
	roles.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return roleService.GetRoleByIDService(c)
	})
	roles.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return roleService.CreateRoleService(c)
	})
	roles.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return roleService.UpdateRoleService(c)
	})
	roles.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return roleService.DeleteRoleService(c)
	})
}
