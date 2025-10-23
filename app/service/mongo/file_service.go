package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	model "go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService struct {
	repo        repository.IFileRepository
	uploadPath  string
}

func NewFileService(repo repository.IFileRepository, uploadPath string) *FileService {
	return &FileService{
		repo:       repo,
		uploadPath: uploadPath,
	}
}

func (s *FileService) UploadFotoService(c *fiber.Ctx) error {
	return s.uploadFile(c, "foto", []string{"image/jpeg", "image/jpg", "image/png"}, 1*1024*1024)
}

func (s *FileService) UploadSertifikatService(c *fiber.Ctx) error {
	return s.uploadFile(c, "sertifikat", []string{"application/pdf"}, 2*1024*1024)
}

func (s *FileService) uploadFile(c *fiber.Ctx, category string, allowedTypes []string, maxSize int64) error {
	// Debug: cek semua form values
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Gagal parsing multipart form",
			"error":   err.Error(),
		})
	}

	// Debug: log semua form values
	fmt.Printf("Form values: %+v\n", form.Value)
	fmt.Printf("Form files: %+v\n", form.File)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File tidak ditemukan dalam request. Pastikan key 'file' ada di form-data",
			"error":   err.Error(),
			"debug":   "Gunakan key 'file' untuk upload file di Postman form-data",
		})
	}

	if fileHeader.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Ukuran file melebihi batas maksimal %d MB", maxSize/(1024*1024)),
		})
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !s.isAllowedType(contentType, allowedTypes) {
		expectedTypes := "JPEG, JPG, PNG"
		if category == "sertifikat" {
			expectedTypes = "PDF"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Tipe file tidak diizinkan. Untuk %s, gunakan format: %s", category, expectedTypes),
			"received_type": contentType,
		})
	}

	alumniID := c.FormValue("alumni_id")
	if alumniID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Alumni ID wajib diisi",
		})
	}

	alumniObjID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Alumni ID tidak valid",
		})
	}

	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadPath, category, newFileName)

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat direktori upload",
			"error":   err.Error(),
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuka file",
			"error":   err.Error(),
		})
	}
	defer file.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat file",
			"error":   err.Error(),
		})
	}
	defer out.Close()

	if _, err := out.ReadFrom(file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menulis file",
			"error":   err.Error(),
		})
	}

	fileModel := &model.File{
		AlumniID:     alumniObjID,
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
		Category:     category,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdFile, err := s.repo.CreateFile(ctx, fileModel)
	if err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan metadata file",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File berhasil diupload",
		"data":    s.toFileResponse(createdFile),
	})
}

func (s *FileService) GetFilesByAlumniIDService(c *fiber.Ctx) error {
	alumniID := c.Params("alumni_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	files, err := s.repo.FindFilesByAlumniID(ctx, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data files dari database",
			"error":   err.Error(),
		})
	}

	var responses []model.FileResponse
	for _, file := range files {
		responses = append(responses, *s.toFileResponse(&file))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data files berhasil diambil dari database",
		"data":    responses,
	})
}

func (s *FileService) DeleteFileService(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file, err := s.repo.FindFileByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data file dari database",
			"error":   err.Error(),
		})
	}

	if file == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File tidak ditemukan",
		})
	}

	if err := os.Remove(file.FilePath); err != nil {
		fmt.Println("Warning: Gagal menghapus file dari storage:", err)
	}

	if err := s.repo.DeleteFile(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menghapus file dari database",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "File berhasil dihapus",
	})
}

func (s *FileService) isAllowedType(contentType string, allowedTypes []string) bool {
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}

func (s *FileService) toFileResponse(file *model.File) *model.FileResponse {
	return &model.FileResponse{
		ID:           file.ID.Hex(),
		AlumniID:     file.AlumniID.Hex(),
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     file.FilePath,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		Category:     file.Category,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}
}
