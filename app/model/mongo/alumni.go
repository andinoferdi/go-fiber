package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	NIM          string              `bson:"nim" json:"nim"`
	Nama         string              `bson:"nama" json:"nama"`
	Jurusan      string              `bson:"jurusan" json:"jurusan"`
	Angkatan     int                 `bson:"angkatan" json:"angkatan"`
	TahunLulus   int                 `bson:"tahun_lulus" json:"tahun_lulus"`
	Email        string              `bson:"email" json:"email"`
	PasswordHash string              `bson:"password_hash" json:"-"`
	NoTelepon    *string             `bson:"no_telepon,omitempty" json:"no_telepon,omitempty"`
	Alamat       *string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	RoleID       primitive.ObjectID  `bson:"role_id" json:"role_id"`
	Role         *Role               `bson:"role,omitempty" json:"role,omitempty"`
	CreatedAt    time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time           `bson:"updated_at" json:"updated_at"`
}

type CreateAlumniRequest struct {
	NIM        string  `bson:"nim" json:"nim" validate:"required"`
	Nama       string  `bson:"nama" json:"nama" validate:"required"`
	Jurusan    string  `bson:"jurusan" json:"jurusan" validate:"required"`
	Angkatan   int     `bson:"angkatan" json:"angkatan" validate:"required"`
	TahunLulus int     `bson:"tahun_lulus" json:"tahun_lulus" validate:"required"`
	Email      string  `bson:"email" json:"email" validate:"required,email"`
	Password   string  `bson:"password" json:"password" validate:"required"`
	NoTelepon  *string `bson:"no_telepon,omitempty" json:"no_telepon"`
	Alamat     *string `bson:"alamat,omitempty" json:"alamat"`
	RoleID     string  `bson:"role_id" json:"role_id" validate:"required"`
}

type UpdateAlumniRequest struct {
	Nama       string  `bson:"nama" json:"nama" validate:"required"`
	Jurusan    string  `bson:"jurusan" json:"jurusan" validate:"required"`
	Angkatan   int     `bson:"angkatan" json:"angkatan" validate:"required"`
	TahunLulus int     `bson:"tahun_lulus" json:"tahun_lulus" validate:"required"`
	Email      string  `bson:"email" json:"email" validate:"required,email"`
	NoTelepon  *string `bson:"no_telepon,omitempty" json:"no_telepon"`
	Alamat     *string `bson:"alamat,omitempty" json:"alamat"`
	RoleID     string  `bson:"role_id" json:"role_id" validate:"required"`
}

