package repository

import (
	"database/sql"
	model "go-fiber/app/model/postgre"
)

func GetAllRoles(db *sql.DB) ([]model.Role, error) {
	query := `
		SELECT id, nama, created_at, updated_at
		FROM roles
		ORDER BY id ASC
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		err := rows.Scan(
			&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func GetRoleByID(db *sql.DB, id int) (*model.Role, error) {
	query := `
		SELECT id, nama, created_at, updated_at
		FROM roles
		WHERE id = $1
	`
	
	role := new(model.Role)
	err := db.QueryRow(query, id).Scan(
		&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return role, nil
}

func GetRoleByName(db *sql.DB, nama string) (*model.Role, error) {
	query := `
		SELECT id, nama, created_at, updated_at
		FROM roles
		WHERE nama = $1
	`
	
	role := new(model.Role)
	err := db.QueryRow(query, nama).Scan(
		&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return role, nil
}

func CreateRole(db *sql.DB, nama string) (*model.Role, error) {
	query := `
		INSERT INTO roles (nama, created_at, updated_at)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, nama, created_at, updated_at
	`
	
	role := new(model.Role)
	err := db.QueryRow(query, nama).Scan(
		&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return role, nil
}

func UpdateRole(db *sql.DB, id int, nama string) (*model.Role, error) {
	query := `
		UPDATE roles 
		SET nama = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, nama, created_at, updated_at
	`
	
	role := new(model.Role)
	err := db.QueryRow(query, nama, id).Scan(
		&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return role, nil
}

func DeleteRole(db *sql.DB, id int) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
