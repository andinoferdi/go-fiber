package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	sortByWhitelist := map[string]bool{"id": true, "alumni_id": true, "nama_perusahaan": true, "posisi_jabatan": true, "bidang_industri": true, "lokasi_kerja": true, "tanggal_mulai_kerja": true, "status_pekerjaan": true, "created_at": true}
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

	pekerjaanList, err := repository.GetAllPekerjaanAlumni(db, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data pekerjaan alumni dari database. Detail: " + err.Error(),
		})
	}

	total, err := repository.CountPekerjaanAlumni(db, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menghitung total data pekerjaan alumni untuk pagination. Detail: " + err.Error(),
		})
	}

	pages := (total + limit - 1) / limit
	if pages == 0 {
		pages = 1
	}

	response := model.PekerjaanAlumniResponse{
		Data: pekerjaanList,
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

func GetPekerjaanAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

	pekerjaan, err := repository.GetPekerjaanAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data pekerjaan alumni dengan ID tersebut tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data pekerjaan alumni dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil diambil dari database.",
		"data":    pekerjaan,
	})
}

func GetPekerjaanAlumniByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	alumniIDStr := c.Params("alumni_id")
	alumniID, err := strconv.Atoi(alumniIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter alumni_id tidak valid. Alumni ID harus berupa angka positif.",
		})
	}

	pekerjaanList, err := repository.GetPekerjaanAlumniByAlumniID(db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data pekerjaan alumni berdasarkan Alumni ID dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berdasarkan Alumni ID berhasil diambil dari database.",
		"data":    pekerjaanList,
	})
}

func CreatePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreatePekerjaanAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
		})
	}

	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" || req.StatusPekerjaan == "" || req.TanggalMulaiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field wajib tidak lengkap. Nama perusahaan, posisi jabatan, bidang industri, lokasi kerja, status pekerjaan, dan tanggal mulai kerja harus diisi.",
		})
	}

	validStatus := map[string]bool{"aktif": true, "selesai": true, "resigned": true}
	if !validStatus[req.StatusPekerjaan] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Status pekerjaan tidak valid. Gunakan 'aktif', 'selesai', atau 'resigned'.",
		})
	}

	tanggalMulaiKerja, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format tanggal mulai kerja tidak valid. Gunakan format YYYY-MM-DD (contoh: 2025-01-15).",
		})
	}

	var tanggalSelesaiKerja *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		parsedTanggalSelesai, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Format tanggal selesai kerja tidak valid. Gunakan format YYYY-MM-DD (contoh: 2025-12-31).",
			})
		}
		tanggalSelesaiKerja = &parsedTanggalSelesai
	}

	pekerjaanRequest := model.CreatePekerjaanAlumniRepositoryRequest{
		AlumniID:            req.AlumniID,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   tanggalMulaiKerja,
		TanggalSelesaiKerja: tanggalSelesaiKerja,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
	}

	pekerjaan, err := repository.CreatePekerjaanAlumni(db, pekerjaanRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menyimpan data pekerjaan alumni ke database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil disimpan ke database.",
		"data":    pekerjaan,
	})
}

func UpdatePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

	var req model.UpdatePekerjaanAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
		})
	}

	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" || req.StatusPekerjaan == "" || req.TanggalMulaiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field wajib tidak lengkap. Nama perusahaan, posisi jabatan, bidang industri, lokasi kerja, status pekerjaan, dan tanggal mulai kerja harus diisi.",
		})
	}

	validStatus := map[string]bool{"aktif": true, "selesai": true, "resigned": true}
	if !validStatus[req.StatusPekerjaan] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Status pekerjaan tidak valid. Gunakan 'aktif', 'selesai', atau 'resigned'.",
		})
	}

	tanggalMulaiKerja, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format tanggal mulai kerja tidak valid. Gunakan format YYYY-MM-DD (contoh: 2025-01-15).",
		})
	}

	var tanggalSelesaiKerja *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		parsedTanggalSelesai, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Format tanggal selesai kerja tidak valid. Gunakan format YYYY-MM-DD (contoh: 2025-12-31).",
			})
		}
		tanggalSelesaiKerja = &parsedTanggalSelesai
	}

	pekerjaanRequest := model.UpdatePekerjaanAlumniRepositoryRequest{
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   tanggalMulaiKerja,
		TanggalSelesaiKerja: tanggalSelesaiKerja,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
	}

	pekerjaan, err := repository.UpdatePekerjaanAlumni(db, id, pekerjaanRequest)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data pekerjaan alumni dengan ID tersebut tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengupdate data pekerjaan alumni di database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil diupdate di database.",
		"data":    pekerjaan,
	})
}

func SoftDeletePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

	alumniID := c.Locals("alumni_id").(int)
	role := c.Locals("role").(string)

	pekerjaan, err := repository.GetPekerjaanAlumniByIDWithDeleted(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data pekerjaan alumni dengan ID tersebut tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data pekerjaan alumni dari database. Detail: " + err.Error(),
		})
	}

	// Cek apakah data sudah di-soft delete
	if pekerjaan.IsDelete != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data pekerjaan alumni sudah di-soft delete sebelumnya.",
		})
	}

	// Validasi permission berdasarkan role
	if role == "admin" {
		err = repository.SoftDeletePekerjaanAlumni(db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error menghapus data pekerjaan alumni dari database. Detail: " + err.Error(),
			})
		}
	} else {
		// User hanya bisa soft delete pekerjaan alumni miliknya sendiri
		if pekerjaan.AlumniID != alumniID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Anda hanya bisa menghapus pekerjaan alumni milik Anda sendiri.",
			})
		}
		
		err = repository.SoftDeletePekerjaanAlumniByAlumniID(db, id, alumniID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error menghapus data pekerjaan alumni dari database. Detail: " + err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil di-soft delete (ditandai sebagai terhapus).",
	})
}

func HardDeletePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

	err = repository.HardDeletePekerjaanAlumni(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menghapus data pekerjaan alumni secara permanen dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil dihapus secara permanen dari database.",
	})
}