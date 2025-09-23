package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"go-fiber/helper"
	"log"

	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var loginReq model.LoginRequest
	
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
	}

	if loginReq.Email == "" || loginReq.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email dan password harus diisi",
		})
	}

	alumni, err := repository.GetAlumniByEmail(db, loginReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[AUTH] Failed login attempt for email: %s from IP: %s", loginReq.Email, c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Email atau password salah",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni: " + err.Error(),
		})
	}

	if !helper.CheckPassword(loginReq.Password, alumni.PasswordHash) {
		log.Printf("[AUTH] Failed login attempt for email: %s from IP: %s", loginReq.Email, c.IP())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Email atau password salah",
		})
	}

	token, err := helper.GenerateToken(*alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal generate token",
			"error":   err.Error(),
		})
	}

	log.Printf("[AUTH] Successful login for alumni: %s (role: %s) from IP: %s", alumni.Email, alumni.Role.Nama, c.IP())
	
	loginResponse := model.LoginResponse{
		Alumni: *alumni,
		Role:   *alumni.Role,
		Token:  token,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    loginResponse,
	})
}

func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	alumniID, ok := c.Locals("alumni_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Alumni information tidak valid",
		})
	}

	alumni, err := repository.GetAlumniByID(db, alumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Alumni tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error database: " + err.Error(),
		})
	}

	profile := model.ProfileResponse{
		AlumniID: alumni.ID,
		NIM:      alumni.NIM,
		Nama:     alumni.Nama,
		Email:    alumni.Email,
		RoleID:   alumni.RoleID,
		RoleName: alumni.Role.Nama,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data":    profile,
	})
}
