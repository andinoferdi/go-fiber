package repository

import (
	"database/sql"
	"go-fiber/app/model"
)

func GetAllRoles(db *sql.DB) ([]model.Role, error) {
	query := `SELECT id, nama, created_at, updated_at FROM roles ORDER BY id`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		err := rows.Scan(&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func GetRoleByID(db *sql.DB, id int) (*model.Role, error) {
	role := new(model.Role)
	query := `SELECT id, nama, created_at, updated_at FROM roles WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func GetRoleByNama(db *sql.DB, nama string) (*model.Role, error) {
	role := new(model.Role)
	query := `SELECT id, nama, created_at, updated_at FROM roles WHERE nama = $1`
	err := db.QueryRow(query, nama).Scan(&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}
