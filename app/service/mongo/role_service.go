package service

import (
	"context"
	model "go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RoleService struct {
	repo repository.IRoleRepository
}

func NewRoleService(repo repository.IRoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAllRolesService(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roles, err := s.repo.FindAllRoles(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data roles dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data roles berhasil diambil dari database.",
		"data":    roles,
	})
}

func (s *RoleService) GetRoleByIDService(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data role dari database. Detail: " + err.Error(),
		})
	}

	if role == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data role dengan ID tersebut tidak ditemukan di database.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data role berhasil diambil dari database.",
		"data":    role,
	})
}

func (s *RoleService) CreateRoleService(c *fiber.Ctx) error {
	var req model.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
		})
	}

	if req.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field nama role wajib diisi.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	role := &model.Role{
		Nama: req.Nama,
	}

	createdRole, err := s.repo.CreateRole(ctx, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menyimpan data role ke database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data role berhasil disimpan ke database.",
		"data":    createdRole,
	})
}

func (s *RoleService) UpdateRoleService(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
		})
	}

	if req.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field nama role wajib diisi.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingRole, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data role dari database. Detail: " + err.Error(),
		})
	}

	if existingRole == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data role dengan ID tersebut tidak ditemukan di database.",
		})
	}

	role := &model.Role{
		Nama: req.Nama,
	}

	updatedRole, err := s.repo.UpdateRole(ctx, id, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengupdate data role di database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data role berhasil diupdate di database.",
		"data":    updatedRole,
	})
}

func (s *RoleService) DeleteRoleService(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data role dari database. Detail: " + err.Error(),
		})
	}

	if role == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data role dengan ID tersebut tidak ditemukan di database.",
		})
	}

	err = s.repo.DeleteRole(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menghapus data role dari database. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data role berhasil dihapus dari database.",
	})
}
