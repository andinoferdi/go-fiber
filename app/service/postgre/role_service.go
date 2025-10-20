package service

import (
	"database/sql"
	model "go-fiber/app/model/postgre"
	repository "go-fiber/app/repository/postgre"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllRolesService(c *fiber.Ctx, db *sql.DB) error {
	roles, err := repository.GetAllRoles(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data roles dari database. Detail: " + err.Error(),
		})
	}

	response := model.GetAllRolesResponse{
		Success: true,
		Message: "Data roles berhasil diambil dari database.",
		Data:    roles,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetRoleByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

	role, err := repository.GetRoleByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data role dengan ID tersebut tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data role dari database. Detail: " + err.Error(),
		})
	}

	response := model.GetRoleByIDResponse{
		Success: true,
		Message: "Data role berhasil diambil dari database.",
		Data:    *role,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateRoleService(c *fiber.Ctx, db *sql.DB) error {
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

	role, err := repository.CreateRole(db, req.Nama)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menyimpan data role ke database. Detail: " + err.Error(),
		})
	}

	response := model.CreateRoleResponse{
		Success: true,
		Message: "Data role berhasil disimpan ke database.",
		Data:    *role,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func UpdateRoleService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

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

	role, err := repository.UpdateRole(db, id, req.Nama)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data role dengan ID tersebut tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengupdate data role di database. Detail: " + err.Error(),
		})
	}

	response := model.UpdateRoleResponse{
		Success: true,
		Message: "Data role berhasil diupdate di database.",
		Data:    *role,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteRoleService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter ID tidak valid. ID harus berupa angka positif.",
		})
	}

	err = repository.DeleteRole(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error menghapus data role dari database. Detail: " + err.Error(),
		})
	}

	response := model.DeleteRoleResponse{
		Success: true,
		Message: "Data role berhasil dihapus dari database.",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}