package middleware

import (
	utilsmongo "go-fiber/utils/mongo"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token akses diperlukan. Tambahkan header 'Authorization: Bearer YOUR_TOKEN'.",
			})
		}

		tokenString := utilsmongo.ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token tidak valid. Gunakan format 'Bearer YOUR_TOKEN'.",
			})
		}

		claims, err := utilsmongo.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau sudah expired. Silakan login ulang untuk mendapatkan token baru.",
			})
		}

		c.Locals("alumni_id", claims.AlumniID)
		c.Locals("email", claims.Email)
		c.Locals("role_id", claims.RoleID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Hanya admin yang dapat mengakses endpoint ini.",
			})
		}
		return c.Next()
	}
}

func UserOrAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" && role != "user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Role tidak valid. Gunakan role 'admin' atau 'user'.",
			})
		}
		return c.Next()
	}
}

