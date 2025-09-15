package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"go-fiber/helper"

	"github.com/gofiber/fiber/v2"
)

// LoginService handles user login authentication
func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var loginReq model.LoginRequest
	
	// Parse request body
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
	}

	// Validate input
	if loginReq.Username == "" || loginReq.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Username dan password harus diisi",
		})
	}

	// Get user from database
	user, passwordHash, err := repository.GetUserByUsername(db, loginReq.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Username atau password salah",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error database: " + err.Error(),
		})
	}

	// Verify password
	if !helper.CheckPassword(loginReq.Password, passwordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Username atau password salah",
		})
	}

	// Generate JWT token
	token, err := helper.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal generate token",
			"error":   err.Error(),
		})
	}

	// Create response
	loginResponse := model.LoginResponse{
		User:  *user,
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    loginResponse,
	})
}

// GetProfileService returns current user profile
func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	// Get user info from middleware context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User information tidak valid",
		})
	}

	// Get full user data from database
	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "User tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error database: " + err.Error(),
		})
	}

	// Create profile response
	profile := model.ProfileResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data":    profile,
	})
}
