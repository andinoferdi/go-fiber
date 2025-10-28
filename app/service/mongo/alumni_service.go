package service

import (
	"context"
	model "go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"
	utilsmongo "go-fiber/utils/mongo"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	repo repository.IAlumniRepository
}

func NewAlumniService(repo repository.IAlumniRepository) *AlumniService {
	return &AlumniService{repo: repo}
}

func (s *AlumniService) GetAllAlumniService(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumniList, err := s.repo.FindAllAlumni(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data alumni dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil dari database.",
		"data":    alumniList,
	})
}

func (s *AlumniService) GetAlumniByIDService(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni, err := s.repo.FindAlumniByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data alumni dari database. Detail: " + err.Error(),
		})
	}

	if alumni == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data alumni dengan ID tersebut tidak ditemukan di database.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil dari database.",
		"data":    alumni,
	})
}

func (s *AlumniService) CreateAlumniService(c *fiber.Ctx) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
		})
	}

	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field wajib tidak lengkap. NIM, nama, jurusan, email, dan password harus diisi.",
		})
	}

	if req.Role != "admin" && req.Role != "user" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Role tidak valid. Gunakan 'admin' atau 'user'.",
		})
	}

	passwordHash, err := utilsmongo.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengenkripsi password. Detail: " + err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni := &model.Alumni{
		NIM:          req.NIM,
		Nama:         req.Nama,
		Jurusan:      req.Jurusan,
		Angkatan:     req.Angkatan,
		TahunLulus:   req.TahunLulus,
		Email:        req.Email,
		PasswordHash: passwordHash,
		NoTelepon:    req.NoTelepon,
		Alamat:       req.Alamat,
		Role:         req.Role,
	}

	createdAlumni, err := s.repo.CreateAlumni(ctx, alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menyimpan data alumni ke database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil disimpan ke database.",
		"data":    createdAlumni,
	})
}

func (s *AlumniService) UpdateAlumniService(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
		})
	}

	if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field wajib tidak lengkap. Nama, jurusan, dan email harus diisi.",
		})
	}

	if req.Role != "admin" && req.Role != "user" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Role tidak valid. Gunakan 'admin' atau 'user'.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingAlumni, err := s.repo.FindAlumniByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data alumni dari database. Detail: " + err.Error(),
		})
	}

	if existingAlumni == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data alumni dengan ID tersebut tidak ditemukan di database.",
		})
	}

	alumni := &model.Alumni{
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
		Role:       req.Role,
	}

	updatedAlumni, err := s.repo.UpdateAlumni(ctx, id, alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengupdate data alumni di database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diupdate di database.",
		"data":    updatedAlumni,
	})
}

func (s *AlumniService) DeleteAlumniService(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni, err := s.repo.FindAlumniByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data alumni dari database. Detail: " + err.Error(),
		})
	}

	if alumni == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data alumni dengan ID tersebut tidak ditemukan di database.",
		})
	}

	err = s.repo.DeleteAlumni(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menghapus data alumni dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil dihapus dari database.",
	})
}

func (s *AlumniService) CheckAlumniService(c *fiber.Ctx) error {
	key := c.Params("key")
	validAPIKey := os.Getenv("API_KEY")
	if validAPIKey == "" {
		validAPIKey = "default-api-key-2024"
	}
	if key != validAPIKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "API key tidak valid. Gunakan key yang benar untuk akses endpoint ini.",
		})
	}

	nim := c.FormValue("nim")
	if nim == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter NIM wajib diisi untuk pengecekan status alumni.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni, err := s.repo.FindAlumniByNIM(ctx, nim)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengecek status alumni di database. Detail: " + err.Error(),
		})
	}

	if alumni == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":  true,
			"message":  "Mahasiswa dengan NIM tersebut bukan alumni.",
			"isAlumni": false,
			"data":     nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"message":  "Data alumni berhasil ditemukan di database.",
		"isAlumni": true,
		"data":     alumni,
	})
}
