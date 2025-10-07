package model

import "time"

type Alumni struct {
	ID           int       `json:"id"`
	NIM          string    `json:"nim"`
	Nama         string    `json:"nama"`
	Jurusan      string    `json:"jurusan"`
	Angkatan     int       `json:"angkatan"`
	TahunLulus   int       `json:"tahun_lulus"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Hidden dari JSON response
	NoTelepon    *string   `json:"no_telepon"`
	Alamat       *string   `json:"alamat"`
	RoleID       int       `json:"role_id"`
	Role         *Role     `json:"role,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateAlumniRequest struct {
	NIM        string  `json:"nim" validate:"required"`
	Nama       string  `json:"nama" validate:"required"`
	Jurusan    string  `json:"jurusan" validate:"required"`
	Angkatan   int     `json:"angkatan" validate:"required"`
	TahunLulus int     `json:"tahun_lulus" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required"`
	NoTelepon  *string `json:"no_telepon"`
	Alamat     *string `json:"alamat"`
	RoleID     int     `json:"role_id" validate:"required"`
}

type UpdateAlumniRequest struct {
	Nama      string  `json:"nama" validate:"required"`
	Jurusan   string  `json:"jurusan" validate:"required"`
	Angkatan  int     `json:"angkatan" validate:"required"`
	TahunLulus int    `json:"tahun_lulus" validate:"required"`
	Email     string  `json:"email" validate:"required,email"`
	NoTelepon *string `json:"no_telepon"`
	Alamat    *string `json:"alamat"`
	RoleID    int     `json:"role_id" validate:"required"`
}
