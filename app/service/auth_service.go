package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"go-fiber/utils"

	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field wajib tidak lengkap. Email dan password harus diisi.",
		})
	}

	alumni, err := repository.GetAlumniByEmail(db, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Login gagal. Email atau password salah.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data alumni dari database. Detail: " + err.Error(),
		})
	}

	if !utils.CheckPassword(req.Password, alumni.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Login gagal. Email atau password salah.",
		})
	}

	token, err := utils.GenerateToken(*alumni, alumni.Role.Nama)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error membuat JWT token. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil. Token JWT telah dibuat.",
		"data": fiber.Map{
			"alumni": alumni,
			"token":  token,
		},
	})
}

func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	alumniID := c.Locals("alumni_id").(int)
	email := c.Locals("email").(string)
	role := c.Locals("role").(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data profile berhasil diambil dari JWT token.",
		"data": fiber.Map{
			"alumni_id": alumniID,
			"email":     email,
			"role":      role,
		},
	})
}