package route

import (
	service "go-fiber/app/service/mongo"
	middleware "go-fiber/middleware/mongo"

	"github.com/gofiber/fiber/v2"
)

func FileRoutes(app *fiber.App, fileService *service.FileService) {
	api := app.Group("/go-fiber-mongo")

	files := api.Group("/files", middleware.AuthRequired())
	files.Get("/", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return fileService.GetAllFilesService(c)
	})
	files.Get("/:id", middleware.UserOrAdmin(), func(c *fiber.Ctx) error {
		return fileService.GetFileByIDService(c)
	})
	files.Post("/upload/foto", middleware.ValidateAlumniAccess(), func(c *fiber.Ctx) error {
		return fileService.UploadFotoService(c)
	})

	files.Post("/upload/sertifikat", middleware.ValidateAlumniAccess(), func(c *fiber.Ctx) error {
		return fileService.UploadSertifikatService(c)
	})

	files.Get("/alumni/:alumni_id", middleware.ValidateAlumniAccess(), func(c *fiber.Ctx) error {
		return fileService.GetFilesByAlumniIDService(c)
	})

	files.Delete("/:id", func(c *fiber.Ctx) error {
		return fileService.DeleteFileService(c)
	})
}
