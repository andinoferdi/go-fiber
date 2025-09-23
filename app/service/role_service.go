package service

import (
	"database/sql"
	"go-fiber/app/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllRolesService(c *fiber.Ctx, db *sql.DB) error {
	roles, err := repository.GetAllRoles(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data roles: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data roles",
		"success": true,
		"data":    roles,
	})
}

func GetRoleByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	role, err := repository.GetRoleByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Role tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data role: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data role",
		"success": true,
		"data":    role,
	})
}
