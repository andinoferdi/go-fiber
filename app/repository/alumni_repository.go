package repository

import (
	"database/sql"
	"fmt"
	"go-fiber/app/model"
	"time"
)

func GetAllAlumni(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	sortByMap := map[string]string{
		"id":           "a.id",
		"nim":          "a.nim",
		"nama":         "a.nama",
		"email":        "a.email",
		"jurusan":      "a.jurusan",
		"angkatan":     "a.angkatan",
		"tahun_lulus":  "a.tahun_lulus",
		"created_at":   "a.created_at",
	}
	
	sortColumn := sortByMap[sortBy]
	if sortColumn == "" {
		sortColumn = "a.id"
	}
	
	query := fmt.Sprintf(`
		SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, 
		       a.no_telepon, a.alamat, a.role_id, a.created_at, a.updated_at,
		       r.id, r.nama, r.created_at, r.updated_at
		FROM alumni a
		LEFT JOIN roles r ON a.role_id = r.id
		WHERE a.nama ILIKE $1 OR a.nim ILIKE $1 OR a.email ILIKE $1 OR a.jurusan ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortColumn, order)
	
	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var alumni model.Alumni
		var role model.Role
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
			&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
			&alumni.NoTelepon, &alumni.Alamat, &alumni.RoleID,
			&alumni.CreatedAt, &alumni.UpdatedAt,
			&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumni.Role = &role
		alumniList = append(alumniList, alumni)
	}

	return alumniList, nil
}

func CountAlumni(db *sql.DB, search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR nim ILIKE $1 OR email ILIKE $1 OR jurusan ILIKE $1`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
	query := `
		SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, 
		       a.no_telepon, a.alamat, a.role_id, a.created_at, a.updated_at,
		       r.id, r.nama, r.created_at, r.updated_at
		FROM alumni a
		LEFT JOIN roles r ON a.role_id = r.id
		WHERE a.id = $1
	`
	
	alumni := new(model.Alumni)
	role := new(model.Role)
	err := db.QueryRow(query, id).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.NoTelepon, &alumni.Alamat, &alumni.RoleID,
		&alumni.CreatedAt, &alumni.UpdatedAt,
		&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	alumni.Role = role
	
	return alumni, nil
}

func CreateAlumni(db *sql.DB, req model.CreateAlumniRequest, passwordHash string) (*model.Alumni, error) {
	query := `
		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, 
		                   password_hash, no_telepon, alamat, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
		          no_telepon, alamat, role_id, created_at, updated_at
	`
	
	now := time.Now()
	alumni := new(model.Alumni)
	err := db.QueryRow(query, 
		req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, passwordHash, req.NoTelepon, req.Alamat, req.RoleID, now, now,
	).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.NoTelepon, &alumni.Alamat, &alumni.RoleID,
		&alumni.CreatedAt, &alumni.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return alumni, nil
}

func UpdateAlumni(db *sql.DB, id int, req model.UpdateAlumniRequest) (*model.Alumni, error) {
	query := `
		UPDATE alumni 
		SET nama = $1, jurusan = $2, angkatan = $3, tahun_lulus = $4, 
		    email = $5, no_telepon = $6, alamat = $7, role_id = $8, updated_at = $9
		WHERE id = $10
		RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
		          no_telepon, alamat, role_id, created_at, updated_at
	`
	
	alumni := new(model.Alumni)
	err := db.QueryRow(query,
		req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, req.RoleID, time.Now(), id,
	).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.NoTelepon, &alumni.Alamat, &alumni.RoleID,
		&alumni.CreatedAt, &alumni.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return alumni, nil
}

func DeleteAlumni(db *sql.DB, id int) error {
	query := `DELETE FROM alumni WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func CheckAlumniByNim(db *sql.DB, nim string) (*model.Alumni, error) {
	query := `
		SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, 
		       a.no_telepon, a.alamat, a.role_id, a.created_at, a.updated_at,
		       r.id, r.nama, r.created_at, r.updated_at
		FROM alumni a
		LEFT JOIN roles r ON a.role_id = r.id
		WHERE a.nim = $1
	`
	
	alumni := new(model.Alumni)
	role := new(model.Role)
	err := db.QueryRow(query, nim).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.NoTelepon, &alumni.Alamat, &alumni.RoleID,
		&alumni.CreatedAt, &alumni.UpdatedAt,
		&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	alumni.Role = role
	
	return alumni, nil
}

func GetAlumniByEmail(db *sql.DB, email string) (*model.Alumni, error) {
	query := `
		SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, 
		       a.password_hash, a.no_telepon, a.alamat, a.role_id, a.created_at, a.updated_at,
		       r.id, r.nama, r.created_at, r.updated_at
		FROM alumni a
		LEFT JOIN roles r ON a.role_id = r.id
		WHERE a.email = $1
	`
	
	alumni := new(model.Alumni)
	role := new(model.Role)
	err := db.QueryRow(query, email).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email,
		&alumni.PasswordHash, &alumni.NoTelepon, &alumni.Alamat, &alumni.RoleID,
		&alumni.CreatedAt, &alumni.UpdatedAt,
		&role.ID, &role.Nama, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	alumni.Role = role
	
	return alumni, nil
}
