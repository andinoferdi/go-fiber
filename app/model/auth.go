package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Alumni Alumni `json:"alumni"`
	Role   Role   `json:"role"`
	Token  string `json:"token"`
}

type JWTClaims struct {
	AlumniID int    `json:"alumni_id"`
	Email    string `json:"email"`
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

type ProfileResponse struct {
	AlumniID int    `json:"alumni_id"`
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}
