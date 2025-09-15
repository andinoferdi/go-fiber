package repository

import (
	"database/sql"
	"go-fiber/app/model"
)

// GetUserByUsername retrieves user by username or email
func GetUserByUsername(db *sql.DB, username string) (*model.User, string, error) {
	var user model.User
	var passwordHash string
	
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at 
		FROM users 
		WHERE username = $1 OR email = $1
	`
	
	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, 
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, "", err
	}
	
	return &user, passwordHash, nil
}

// GetUserByID retrieves user by ID
func GetUserByID(db *sql.DB, id int) (*model.User, error) {
	var user model.User
	
	query := `
		SELECT id, username, email, role, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`
	
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, 
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// CreateUser creates a new user (for future use)
func CreateUser(db *sql.DB, username, email, passwordHash, role string) (*model.User, error) {
	var user model.User
	
	query := `
		INSERT INTO users (username, email, password_hash, role, updated_at) 
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) 
		RETURNING id, username, email, role, created_at, updated_at
	`
	
	err := db.QueryRow(query, username, email, passwordHash, role).Scan(
		&user.ID, &user.Username, &user.Email, 
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}
