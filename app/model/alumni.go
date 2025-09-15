package model

import "time"

type Alumni struct {
	ID          int       `json:"id" db:"id"`
	NIM         string    `json:"nim" db:"nim"`
	Nama        string    `json:"nama" db:"nama"`
	Jurusan     string    `json:"jurusan" db:"jurusan"`
	Angkatan    int       `json:"angkatan" db:"angkatan"`
	TahunLulus  int       `json:"tahun_lulus" db:"tahun_lulus"`
	Email       string    `json:"email" db:"email"`
	NoTelepon   *string   `json:"no_telepon" db:"no_telepon"`
	Alamat      *string   `json:"alamat" db:"alamat"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
