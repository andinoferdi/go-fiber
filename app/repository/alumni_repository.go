package repository

import (
	"database/sql"
	"fmt"
	"go-fiber/app/model"
	"time"
)

func GetAllAlumni(db *sql.DB) ([]model.Alumni, error) {
	query := `SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.password_hash, a.role_id, a.no_telepon, a.alamat, a.created_at, a.updated_at, r.id, r.nama, r.created_at, r.updated_at FROM alumni a LEFT JOIN roles r ON a.role_id = r.id ORDER BY a.created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumni []model.Alumni
	for rows.Next() {
		var a model.Alumni
		var role model.Role
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.PasswordHash, &a.RoleID, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt, &role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		a.Role = &role
		alumni = append(alumni, a)
	}
	return alumni, nil
}

func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	var role model.Role
	query := `SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.password_hash, a.role_id, a.no_telepon, a.alamat, a.created_at, a.updated_at, r.id, r.nama, r.created_at, r.updated_at FROM alumni a LEFT JOIN roles r ON a.role_id = r.id WHERE a.id = $1`
	err := db.QueryRow(query, id).Scan(&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.PasswordHash, &alumni.RoleID, &alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt, &role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	alumni.Role = &role
	return alumni, nil
}

func GetAlumniByEmail(db *sql.DB, email string) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	var role model.Role
	query := `SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.password_hash, a.role_id, a.no_telepon, a.alamat, a.created_at, a.updated_at, r.id, r.nama, r.created_at, r.updated_at FROM alumni a LEFT JOIN roles r ON a.role_id = r.id WHERE a.email = $1`
	err := db.QueryRow(query, email).Scan(&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.PasswordHash, &alumni.RoleID, &alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt, &role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	alumni.Role = &role
	return alumni, nil
}

func CreateAlumni(db *sql.DB, alumni *model.Alumni) error {
	query := `INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, password_hash, role_id, no_telepon, alamat, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
			  RETURNING id, created_at`
	alumni.UpdatedAt = time.Now()
	err := db.QueryRow(query, alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.PasswordHash, alumni.RoleID, alumni.NoTelepon, alumni.Alamat, alumni.UpdatedAt).Scan(&alumni.ID, &alumni.CreatedAt)
	return err
}

func UpdateAlumni(db *sql.DB, id int, alumni *model.Alumni) error {
	query := `UPDATE alumni SET nim = $1, nama = $2, jurusan = $3, angkatan = $4, tahun_lulus = $5, email = $6, password_hash = $7, role_id = $8, no_telepon = $9, alamat = $10, updated_at = $11 WHERE id = $12`
	alumni.UpdatedAt = time.Now()
	_, err := db.Exec(query, alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.PasswordHash, alumni.RoleID, alumni.NoTelepon, alumni.Alamat, alumni.UpdatedAt, id)
	return err
}

func DeleteAlumni(db *sql.DB, id int) error {
	query := `DELETE FROM alumni WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func GetAlumniWithPagination(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	query := fmt.Sprintf(`
		SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.password_hash, a.role_id, a.no_telepon, a.alamat, a.created_at, a.updated_at, r.id, r.nama, r.created_at, r.updated_at
		FROM alumni a
		LEFT JOIN roles r ON a.role_id = r.id
		WHERE a.nama ILIKE $1 OR a.email ILIKE $1 OR a.nim ILIKE $1 OR a.jurusan ILIKE $1
		ORDER BY a.%s %s
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
		var role model.Role
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.PasswordHash, &a.RoleID, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt, &role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		a.Role = &role
		alumni = append(alumni, a)
	}
	return alumni, nil
}

func CountAlumni(db *sql.DB, search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR email ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}
