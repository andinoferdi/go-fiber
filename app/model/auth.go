package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents user data in the database
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents login response with user data and JWT token
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// JWTClaims represents JWT token claims structure
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// ProfileResponse represents user profile response
type ProfileResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
