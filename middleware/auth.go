package middleware

import (
	"go-fiber/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token akses diperlukan",
				"error":   "Authorization header missing",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token tidak valid",
				"error":   "Invalid authorization header format",
			})
		}

		claims, err := helper.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau expired",
				"error":   err.Error(),
			})
		}

		c.Locals("alumni_id", claims.AlumniID)
		c.Locals("email", claims.Email)
		c.Locals("role_id", claims.RoleID)
		c.Locals("role_name", claims.RoleName)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName := c.Locals("role_name")
		if roleName == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token akses diperlukan",
				"error":   "Role information not found in context",
			})
		}

		userRole, ok := roleName.(string)
		if !ok || userRole != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Hanya admin yang diizinkan",
				"error":   "Insufficient privileges",
			})
		}

		return c.Next()
	}
}

func UserOrAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName := c.Locals("role_name")
		if roleName == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token akses diperlukan",
				"error":   "Role information not found in context",
			})
		}

		userRole, ok := roleName.(string)
		if !ok || (userRole != "admin" && userRole != "user") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Role tidak valid",
				"error":   "Invalid role",
			})
		}

		return c.Next()
	}
}
