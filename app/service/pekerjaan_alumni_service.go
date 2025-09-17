package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetAllPekerjaanAlumniService - Service untuk mengambil semua data pekerjaan alumni dengan pagination, sorting, dan search
func GetAllPekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	// Calculate offset
	offset := (page - 1) * limit

	// Validasi input
	sortByWhitelist := map[string]bool{
		"id": true, "alumni_id": true, "nama_perusahaan": true, "posisi_jabatan": true, 
		"bidang_industri": true, "lokasi_kerja": true, "tanggal_mulai_kerja": true, 
		"status_pekerjaan": true, "created_at": true,
	}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}

	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// Ambil data dari repository dengan pagination
	pekerjaan, err := repository.GetPekerjaanAlumniWithPagination(db, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	// Hitung total data
	total, err := repository.CountPekerjaanAlumni(db, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghitung total data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	// Buat response dengan pagination
	response := model.PekerjaanAlumniResponse{
		Data: pekerjaan,
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
		"message": "Berhasil mengambil data pekerjaan",
		"success": true,
		"data":    response.Data,
		"meta":    response.Meta,
	})
}

// GetPekerjaanAlumniByIDService - Service untuk mengambil data pekerjaan berdasarkan ID
func GetPekerjaanAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	pekerjaan, err := repository.GetPekerjaanAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Data pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data pekerjaan",
		"success": true,
		"data":    pekerjaan,
	})
}

// GetPekerjaanByAlumniIDService - Service untuk mengambil semua pekerjaan berdasarkan alumni ID
func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	alumniIDStr := c.Params("alumni_id")
	alumniID, err := strconv.Atoi(alumniIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Alumni ID tidak valid",
			"success": false,
		})
	}

	// Cek apakah alumni ada
	_, err = repository.GetAlumniByID(db, alumniID)
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

	pekerjaan, err := repository.GetPekerjaanByAlumniID(db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data pekerjaan alumni",
		"success": true,
		"data":    pekerjaan,
	})
}

// CreatePekerjaanAlumniService - Service untuk menambah pekerjaan baru
func CreatePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	var pekerjaan model.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Validasi field wajib
	if pekerjaan.AlumniID == 0 || pekerjaan.NamaPerusahaan == "" || pekerjaan.PosisiJabatan == "" || pekerjaan.BidangIndustri == "" || pekerjaan.LokasiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Field Alumni ID, Nama Perusahaan, Posisi Jabatan, Bidang Industri, dan Lokasi Kerja wajib diisi",
			"success": false,
		})
	}

	// Cek apakah alumni ada
	_, err := repository.GetAlumniByID(db, pekerjaan.AlumniID)
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

	// Set default status jika kosong
	if pekerjaan.StatusPekerjaan == "" {
		pekerjaan.StatusPekerjaan = "aktif"
	}

	// Validasi status pekerjaan
	if pekerjaan.StatusPekerjaan != "aktif" && pekerjaan.StatusPekerjaan != "selesai" && pekerjaan.StatusPekerjaan != "resigned" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Status pekerjaan harus 'aktif', 'selesai', atau 'resigned'",
			"success": false,
		})
	}

	if err := repository.CreatePekerjaanAlumni(db, &pekerjaan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambah pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil menambah pekerjaan",
		"success": true,
		"data":    pekerjaan,
	})
}

// UpdatePekerjaanAlumniService - Service untuk mengupdate data pekerjaan
func UpdatePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	var pekerjaan model.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Validasi field wajib
	if pekerjaan.AlumniID == 0 || pekerjaan.NamaPerusahaan == "" || pekerjaan.PosisiJabatan == "" || pekerjaan.BidangIndustri == "" || pekerjaan.LokasiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Field Alumni ID, Nama Perusahaan, Posisi Jabatan, Bidang Industri, dan Lokasi Kerja wajib diisi",
			"success": false,
		})
	}

	// Cek apakah pekerjaan ada
	_, err = repository.GetPekerjaanAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Data pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengecek data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	// Cek apakah alumni ada
	_, err = repository.GetAlumniByID(db, pekerjaan.AlumniID)
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

	// Validasi status pekerjaan
	if pekerjaan.StatusPekerjaan != "aktif" && pekerjaan.StatusPekerjaan != "selesai" && pekerjaan.StatusPekerjaan != "resigned" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Status pekerjaan harus 'aktif', 'selesai', atau 'resigned'",
			"success": false,
		})
	}

	if err := repository.UpdatePekerjaanAlumni(db, id, &pekerjaan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	pekerjaan.ID = id
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengupdate pekerjaan",
		"success": true,
		"data":    pekerjaan,
	})
}

// DeletePekerjaanAlumniService - Service untuk menghapus data pekerjaan
func DeletePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	// Cek apakah pekerjaan ada
	_, err = repository.GetPekerjaanAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Data pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengecek data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if err := repository.DeletePekerjaanAlumni(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menghapus pekerjaan",
		"success": true,
	})
}
