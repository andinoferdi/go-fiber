package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
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

	sortByWhitelist := map[string]bool{
		"id": true, "nim": true, "nama": true, "jurusan": true, 
		"angkatan": true, "tahun_lulus": true, "email": true, "created_at": true,
	}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}

	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	alumni, err := repository.GetAlumniWithPagination(db, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data alumni: " + err.Error(),
			"success": false,
		})
	}

	total, err := repository.CountAlumni(db, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghitung total data alumni: " + err.Error(),
			"success": false,
		})
	}

	response := model.AlumniResponse{
		Data: alumni,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data alumni",
		"success": true,
		"data":    response.Data,
		"meta":    response.Meta,
	})
}

func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	alumni, err := repository.GetAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alumni tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data alumni",
		"success": true,
		"data":    alumni,
	})
}

func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	if alumni.NIM == "" || alumni.Nama == "" || alumni.Jurusan == "" || alumni.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Field NIM, Nama, Jurusan, dan Email wajib diisi",
			"success": false,
		})
	}

	if err := repository.CreateAlumni(db, &alumni); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambah alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil menambah alumni",
		"success": true,
		"data":    alumni,
	})
}

func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	if alumni.NIM == "" || alumni.Nama == "" || alumni.Jurusan == "" || alumni.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Field NIM, Nama, Jurusan, dan Email wajib diisi",
			"success": false,
		})
	}

	_, err = repository.GetAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alumni tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengecek data alumni: " + err.Error(),
			"success": false,
		})
	}

	if err := repository.UpdateAlumni(db, id, &alumni); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate alumni: " + err.Error(),
			"success": false,
		})
	}

	alumni.ID = id
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengupdate alumni",
		"success": true,
		"data":    alumni,
	})
}

func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	_, err = repository.GetAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alumni tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengecek data alumni: " + err.Error(),
			"success": false,
		})
	}

	if err := repository.DeleteAlumni(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menghapus alumni",
		"success": true,
	})
}
