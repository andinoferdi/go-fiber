package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func ValidateAlumniAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		alumniIDFromToken := c.Locals("alumni_id").(string)
		alumniIDFromParam := c.Params("alumni_id")
		alumniIDFromForm := c.FormValue("alumni_id")

		if role == "admin" {
			return c.Next()
		}

		if role == "user" {
			// Cek dari URL parameter (untuk GET /files/alumni/:alumni_id)
			if alumniIDFromParam != "" && alumniIDFromParam != alumniIDFromToken {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"message": "Anda hanya bisa menambahkan foto dan sertifikat sendiri",
				})
			}

			// Cek dari form data (untuk upload file)
			if alumniIDFromForm != "" && alumniIDFromForm != alumniIDFromToken {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"message": "Anda hanya bisa menambahkan foto dan sertifikat diri anda sendiri",
				})
			}
		}

		return c.Next()
	}
}
