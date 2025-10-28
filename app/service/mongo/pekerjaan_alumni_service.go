package service

import (
	"context"
	model "go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanAlumniService struct {
	repo repository.IPekerjaanAlumniRepository
}

func NewPekerjaanAlumniService(repo repository.IPekerjaanAlumniRepository) *PekerjaanAlumniService {
	return &PekerjaanAlumniService{repo: repo}
}

func (s *PekerjaanAlumniService) GetAllPekerjaanAlumniService(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaanList, err := s.repo.FindAllPekerjaanAlumni(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data pekerjaan alumni dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil diambil dari database.",
		"data":    pekerjaanList,
	})
}

func (s *PekerjaanAlumniService) GetPekerjaanAlumniByIDService(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan, err := s.repo.FindPekerjaanAlumniByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data pekerjaan alumni dari database. Detail: " + err.Error(),
		})
	}

	if pekerjaan == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data pekerjaan alumni dengan ID tersebut tidak ditemukan di database.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil diambil dari database.",
		"data":    pekerjaan,
	})
}

func (s *PekerjaanAlumniService) GetPekerjaanAlumniByAlumniIDService(c *fiber.Ctx) error {
	alumniID := c.Params("alumni_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaanList, err := s.repo.FindPekerjaanAlumniByAlumniID(ctx, alumniID)
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

func (s *PekerjaanAlumniService) CreatePekerjaanAlumniService(c *fiber.Ctx) error {
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

	// Validasi AlumniInfo
	if req.AlumniInfo.NIM == "" || req.AlumniInfo.Nama == "" || req.AlumniInfo.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "AlumniInfo tidak lengkap. NIM, Nama, dan Email harus diisi.",
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	pekerjaan := &model.PekerjaanAlumni{
		AlumniInfo:          req.AlumniInfo,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   tanggalMulaiKerja,
		TanggalSelesaiKerja: tanggalSelesaiKerja,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	createdPekerjaan, err := s.repo.CreatePekerjaanAlumni(ctx, pekerjaan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menyimpan data pekerjaan alumni ke database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil disimpan ke database.",
		"data":    createdPekerjaan,
	})
}

func (s *PekerjaanAlumniService) UpdatePekerjaanAlumniService(c *fiber.Ctx) error {
	id := c.Params("id")
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingPekerjaan, err := s.repo.FindPekerjaanAlumniByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data pekerjaan alumni dari database. Detail: " + err.Error(),
		})
	}

	if existingPekerjaan == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data pekerjaan alumni dengan ID tersebut tidak ditemukan di database.",
		})
	}

	pekerjaan := &model.PekerjaanAlumni{
		AlumniInfo:          existingPekerjaan.AlumniInfo,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   tanggalMulaiKerja,
		TanggalSelesaiKerja: tanggalSelesaiKerja,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		UpdatedAt:           time.Now(),
	}

	updatedPekerjaan, err := s.repo.UpdatePekerjaanAlumni(ctx, id, pekerjaan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengupdate data pekerjaan alumni di database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil diupdate di database.",
		"data":    updatedPekerjaan,
	})
}

func (s *PekerjaanAlumniService) DeletePekerjaanAlumniService(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan, err := s.repo.FindPekerjaanAlumniByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data pekerjaan alumni dari database. Detail: " + err.Error(),
		})
	}

	if pekerjaan == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data pekerjaan alumni dengan ID tersebut tidak ditemukan di database.",
		})
	}

	err = s.repo.DeletePekerjaanAlumni(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menghapus data pekerjaan alumni dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil dihapus dari database.",
	})
}
