package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"go-fiber/utils"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	sortByWhitelist := map[string]bool{"id": true, "nim": true, "nama": true, "email": true, "jurusan": true, "angkatan": true, "tahun_lulus": true, "created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}

	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	alumniList, err := repository.GetAllAlumni(db, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data alumni dari database. Detail: " + err.Error(),
		})
	}

	total, err := repository.CountAlumni(db, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menghitung total data alumni untuk pagination. Detail: " + err.Error(),
		})
	}

	pages := (total + limit - 1) / limit
	if pages == 0 {
		pages = 1
	}

	response := model.AlumniResponse{
		Data: alumniList,
		Meta: model.MetaInfo{
			Page:    page,
			Limit:   limit,
			Total:   total,
			Pages:   pages,
			SortBy:  sortBy,
			Order:   order,
			Search:  search,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

	alumni, err := repository.GetAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data alumni dengan ID tersebut tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data alumni dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil dari database.",
		"data":    alumni,
	})
}

func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

	if req.RoleID != 1 && req.RoleID != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Role ID tidak valid. Gunakan 1 untuk admin atau 2 untuk user.",
		})
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengenkripsi password. Detail: " + err.Error(),
		})
	}

	alumni, err := repository.CreateAlumni(db, req, passwordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menyimpan data alumni ke database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil disimpan ke database.",
		"data":    alumni,
	})
}

func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

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

	if req.RoleID != 1 && req.RoleID != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Role ID tidak valid. Gunakan 1 untuk admin atau 2 untuk user.",
		})
	}

	alumni, err := repository.UpdateAlumni(db, id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data alumni dengan ID tersebut tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengupdate data alumni di database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diupdate di database.",
		"data":    alumni,
	})
}

func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

	err = repository.DeleteAlumni(db, id)
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

func CheckAlumniService(c *fiber.Ctx, db *sql.DB) error {
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
	
	alumni, err := repository.CheckAlumniByNim(db, nim)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success":  true,
				"message":  "Mahasiswa dengan NIM tersebut bukan alumni.",
				"isAlumni": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengecek status alumni di database. Detail: " + err.Error(),
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"message":  "Data alumni berhasil ditemukan di database.",
		"isAlumni": true,
		"data":     alumni,
	})
}