package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

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

	pekerjaan, err := repository.GetPekerjaanAlumniWithPagination(db, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	total, err := repository.CountPekerjaanAlumni(db, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghitung total data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

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

func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	alumniIDStr := c.Params("alumni_id")
	alumniID, err := strconv.Atoi(alumniIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Alumni ID tidak valid",
			"success": false,
		})
	}

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

	if pekerjaan.StatusPekerjaan == "" {
		pekerjaan.StatusPekerjaan = "aktif"
	}

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

	if pekerjaan.AlumniID == 0 || pekerjaan.NamaPerusahaan == "" || pekerjaan.PosisiJabatan == "" || pekerjaan.BidangIndustri == "" || pekerjaan.LokasiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Field Alumni ID, Nama Perusahaan, Posisi Jabatan, Bidang Industri, dan Lokasi Kerja wajib diisi",
			"success": false,
		})
	}

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

func DeletePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

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


func SoftDeletePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	userID, ok := c.Locals("alumni_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Alumni information tidak valid",
			"success": false,
		})
	}

	userRole, ok := c.Locals("role_name").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Role information tidak valid",
			"success": false,
		})
	}

	pekerjaan, err := repository.GetPekerjaanAlumniByIDForSoftDelete(db, id)
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

	if pekerjaan.IsDelete != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data pekerjaan sudah dihapus",
			"success": false,
		})
	}

	if userRole == "admin" {
		if err := repository.SoftDeletePekerjaanAlumni(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menghapus pekerjaan: " + err.Error(),
				"success": false,
			})
		}
	} else if userRole == "user" {
		// User hanya bisa menghapus pekerjaan yang terkait dengan alumni mereka
		if pekerjaan.AlumniID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Anda hanya bisa menghapus pekerjaan yang berelasi dengan akun Anda",
				"success": false,
			})
		}

		if err := repository.SoftDeletePekerjaanAlumni(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menghapus pekerjaan: " + err.Error(),
				"success": false,
			})
		}
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Role tidak valid",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menghapus pekerjaan",
		"success": true,
	})
}
