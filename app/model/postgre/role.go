package model

import "time"

type Role struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRoleRequest struct {
	Nama string `json:"nama" validate:"required"`
}

type UpdateRoleRequest struct {
	Nama string `json:"nama" validate:"required"`
}

type GetAllRolesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []Role `json:"data"`
}

type GetRoleByIDResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type CreateRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type UpdateRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type DeleteRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

