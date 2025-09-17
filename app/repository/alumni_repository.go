package repository

import (
	"database/sql"
	"fmt"
	"go-fiber/app/model"
	"time"
)

// GetAllAlumni - Ambil semua data alumni
func GetAllAlumni(db *sql.DB) ([]model.Alumni, error) {
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumni []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		alumni = append(alumni, a)
	}
	return alumni, nil
}

// GetAlumniByID - Ambil data alumni berdasarkan ID
func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

// CreateAlumni - Tambah alumni baru
func CreateAlumni(db *sql.DB, alumni *model.Alumni) error {
	query := `INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			  RETURNING id, created_at`
	alumni.UpdatedAt = time.Now()
	err := db.QueryRow(query, alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, alumni.UpdatedAt).Scan(&alumni.ID, &alumni.CreatedAt)
	return err
}

// UpdateAlumni - Update data alumni
func UpdateAlumni(db *sql.DB, id int, alumni *model.Alumni) error {
	query := `UPDATE alumni SET nim = $1, nama = $2, jurusan = $3, angkatan = $4, tahun_lulus = $5, email = $6, no_telepon = $7, alamat = $8, updated_at = $9 WHERE id = $10`
	alumni.UpdatedAt = time.Now()
	_, err := db.Exec(query, alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, alumni.UpdatedAt, id)
	return err
}

// DeleteAlumni - Hapus data alumni
func DeleteAlumni(db *sql.DB, id int) error {
	query := `DELETE FROM alumni WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// GetAlumniWithPagination - Ambil data alumni dengan pagination, sorting, dan search
func GetAlumniWithPagination(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	query := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE nama ILIKE $1 OR email ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumni []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		alumni = append(alumni, a)
	}
	return alumni, nil
}

// CountAlumni - Hitung total data alumni untuk pagination
func CountAlumni(db *sql.DB, search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR email ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}
