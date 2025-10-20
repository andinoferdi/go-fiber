package repository

import (
	"database/sql"
	"fmt"
	model "go-fiber/app/model/postgre"
	"time"
)

func GetAllPekerjaanAlumni(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at
		FROM pekerjaan_alumni 
		WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1)
		AND is_delete IS NULL
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)
	
	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
			&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
			&pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja,
			&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.IsDelete,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, nil
}

func CountPekerjaanAlumni(db *sql.DB, search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1) AND is_delete IS NULL`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func GetPekerjaanAlumniByID(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at
		FROM pekerjaan_alumni 
		WHERE id = $1 AND is_delete IS NULL
	`
	
	pekerjaan := new(model.PekerjaanAlumni)
	err := db.QueryRow(query, id).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja,
		&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.IsDelete,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return pekerjaan, nil
}

func GetPekerjaanAlumniByAlumniID(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at
		FROM pekerjaan_alumni 
		WHERE alumni_id = $1 AND is_delete IS NULL
		ORDER BY tanggal_mulai_kerja DESC
	`
	
	rows, err := db.Query(query, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
			&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
			&pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja,
			&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.IsDelete,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, nil
}

func CreatePekerjaanAlumni(db *sql.DB, req model.CreatePekerjaanAlumniRepositoryRequest) (*model.PekerjaanAlumni, error) {
	query := `
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, 
		                             bidang_industri, lokasi_kerja, gaji_range,
		                             tanggal_mulai_kerja, tanggal_selesai_kerja,
		                             status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		          lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		          status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at
	`
	
	now := time.Now()
	pekerjaan := new(model.PekerjaanAlumni)
	err := db.QueryRow(query,
		req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
		req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja,
		req.StatusPekerjaan, req.DeskripsiPekerjaan, nil, now, now,
	).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja,
		&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.IsDelete,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return pekerjaan, nil
}

func UpdatePekerjaanAlumni(db *sql.DB, id int, req model.UpdatePekerjaanAlumniRepositoryRequest) (*model.PekerjaanAlumni, error) {
	query := `
		UPDATE pekerjaan_alumni 
		SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3,
		    lokasi_kerja = $4, gaji_range = $5, tanggal_mulai_kerja = $6,
		    tanggal_selesai_kerja = $7, status_pekerjaan = $8, deskripsi_pekerjaan = $9,
		    updated_at = $10
		WHERE id = $11 AND is_delete IS NULL
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		          lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		          status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at
	`
	
	pekerjaan := new(model.PekerjaanAlumni)
	err := db.QueryRow(query,
		req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja,
		req.StatusPekerjaan, req.DeskripsiPekerjaan, time.Now(), id,
	).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja,
		&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.IsDelete,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return pekerjaan, nil
}

func SoftDeletePekerjaanAlumni(db *sql.DB, id int) error {
	query := `UPDATE pekerjaan_alumni SET is_delete = $1 WHERE id = $2`
	_, err := db.Exec(query, time.Now(), id)
	return err
}

func SoftDeletePekerjaanAlumniByAlumniID(db *sql.DB, id int, alumniID int) error {
	query := `UPDATE pekerjaan_alumni SET is_delete = $1 WHERE id = $2 AND alumni_id = $3`
	_, err := db.Exec(query, time.Now(), id, alumniID)
	return err
}

func GetPekerjaanAlumniByIDWithDeleted(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at
		FROM pekerjaan_alumni 
		WHERE id = $1
	`
	
	pekerjaan := new(model.PekerjaanAlumni)
	err := db.QueryRow(query, id).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja,
		&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.IsDelete,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return pekerjaan, nil
}

func HardDeletePekerjaanAlumni(db *sql.DB, id int) error {
	query := `DELETE FROM pekerjaan_alumni WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func GetSoftDeletedPekerjaanAlumni(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at
		FROM pekerjaan_alumni 
		WHERE alumni_id = $1 AND is_delete IS NOT NULL
		ORDER BY is_delete DESC
	`
	
	rows, err := db.Query(query, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
			&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
			&pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja,
			&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.IsDelete,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, nil
}

func GetAllSoftDeletedPekerjaanAlumni(db *sql.DB) ([]model.PekerjaanAlumni, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at, updated_at
		FROM pekerjaan_alumni 
		WHERE is_delete IS NOT NULL
		ORDER BY is_delete DESC
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan,
			&pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja,
			&pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja,
			&pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.IsDelete,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, nil
}

func RestorePekerjaanAlumni(db *sql.DB, id int) error {
	query := `UPDATE pekerjaan_alumni SET is_delete = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}