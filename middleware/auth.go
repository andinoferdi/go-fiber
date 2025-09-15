package middleware

import (
	"go-fiber/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthRequired middleware memverifikasi JWT token dan menyimpan user info di context
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token akses diperlukan",
				"error":   "Authorization header missing",
			})
		}

		// Extract token dari "Bearer TOKEN"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token tidak valid",
				"error":   "Invalid authorization header format",
			})
		}

		// Validasi token
		claims, err := helper.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau expired",
				"error":   err.Error(),
			})
		}

		// Simpan informasi user di context untuk handler berikutnya
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// AdminOnly middleware memastikan user memiliki role admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token akses diperlukan",
				"error":   "User information not found in context",
			})
		}

		userRole, ok := role.(string)
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

// UserOrAdmin middleware memungkinkan akses untuk user biasa atau admin
func UserOrAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token akses diperlukan",
				"error":   "User information not found in context",
			})
		}

		userRole, ok := role.(string)
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
