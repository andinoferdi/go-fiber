package service

import (
	"context"
	"time"

	modelMongo "go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"
	utilsmongo "go-fiber/utils/mongo"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	alumniRepo repository.IAlumniRepository
}

func NewAuthService(alumniRepo repository.IAlumniRepository) *AuthService {
	return &AuthService{alumniRepo: alumniRepo}
}

func (s *AuthService) LoginService(c *fiber.Ctx) error {
	var req modelMongo.LoginRequest
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni, err := s.alumniRepo.FindAlumniByEmail(ctx, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data alumni dari database. Detail: " + err.Error(),
		})
	}

	if alumni == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Login gagal. Email atau password salah.",
		})
	}

	if !utilsmongo.CheckPassword(req.Password, alumni.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Login gagal. Email atau password salah.",
		})
	}

	alumniToken := utilsmongo.AlumniToken{
		ID:     alumni.ID.Hex(),
		Email:  alumni.Email,
		Role:   alumni.Role,
	}

	token, err := utilsmongo.GenerateToken(alumniToken)
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

func (s *AuthService) GetProfileService(c *fiber.Ctx) error {
	alumniID := c.Locals("alumni_id")
	email := c.Locals("email")
	role := c.Locals("role")

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

