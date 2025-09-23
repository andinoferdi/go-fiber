package repository

import (
	"database/sql"
	"fmt"
	"go-fiber/app/model"
	"time"
)

func GetAllPekerjaanAlumni(db *sql.DB) ([]model.PekerjaanAlumni, error) {
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete FROM pekerjaan_alumni WHERE is_delete IS NULL ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var tanggalMulai, tanggalSelesai sql.NullTime
		var isDelete sql.NullTime
		
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &tanggalMulai, &tanggalSelesai, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &isDelete)
		if err != nil {
			return nil, err
		}
		
		if tanggalMulai.Valid {
			p.TanggalMulaiKerja = model.CustomDate{Time: tanggalMulai.Time}
		}
		if tanggalSelesai.Valid {
			p.TanggalSelesaiKerja = &model.CustomDate{Time: tanggalSelesai.Time}
		}
		if isDelete.Valid {
			p.IsDelete = &isDelete.Time
		}
		
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func GetPekerjaanAlumniByID(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	pekerjaan := new(model.PekerjaanAlumni)
	var tanggalMulai, tanggalSelesai sql.NullTime
	var isDelete sql.NullTime
	
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete FROM pekerjaan_alumni WHERE id = $1 AND is_delete IS NULL`
	err := db.QueryRow(query, id).Scan(&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan, &pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja, &pekerjaan.GajiRange, &tanggalMulai, &tanggalSelesai, &pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt, &isDelete)
	if err != nil {
		return nil, err
	}
	
	if tanggalMulai.Valid {
		pekerjaan.TanggalMulaiKerja = model.CustomDate{Time: tanggalMulai.Time}
	}
	if tanggalSelesai.Valid {
		pekerjaan.TanggalSelesaiKerja = &model.CustomDate{Time: tanggalSelesai.Time}
	}
	if isDelete.Valid {
		pekerjaan.IsDelete = &isDelete.Time
	}
	
	return pekerjaan, nil
}

func GetPekerjaanByAlumniID(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete FROM pekerjaan_alumni WHERE alumni_id = $1 AND is_delete IS NULL ORDER BY tanggal_mulai_kerja DESC`
	rows, err := db.Query(query, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var tanggalMulai, tanggalSelesai sql.NullTime
		var isDelete sql.NullTime
		
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &tanggalMulai, &tanggalSelesai, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &isDelete)
		if err != nil {
			return nil, err
		}
		
		if tanggalMulai.Valid {
			p.TanggalMulaiKerja = model.CustomDate{Time: tanggalMulai.Time}
		}
		if tanggalSelesai.Valid {
			p.TanggalSelesaiKerja = &model.CustomDate{Time: tanggalSelesai.Time}
		}
		if isDelete.Valid {
			p.IsDelete = &isDelete.Time
		}
		
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func CreatePekerjaanAlumni(db *sql.DB, pekerjaan *model.PekerjaanAlumni) error {
	query := `INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
			  RETURNING id, created_at`
	pekerjaan.UpdatedAt = time.Now()
	
	var tanggalSelesai interface{}
	if pekerjaan.TanggalSelesaiKerja != nil && !pekerjaan.TanggalSelesaiKerja.Time.IsZero() {
		tanggalSelesai = pekerjaan.TanggalSelesaiKerja.Time
	} else {
		tanggalSelesai = nil
	}
	
	err := db.QueryRow(query, 
		pekerjaan.AlumniID, 
		pekerjaan.NamaPerusahaan, 
		pekerjaan.PosisiJabatan, 
		pekerjaan.BidangIndustri, 
		pekerjaan.LokasiKerja, 
		pekerjaan.GajiRange, 
		pekerjaan.TanggalMulaiKerja.Time, 
		tanggalSelesai, 
		pekerjaan.StatusPekerjaan, 
		pekerjaan.DeskripsiPekerjaan, 
		pekerjaan.UpdatedAt).Scan(&pekerjaan.ID, &pekerjaan.CreatedAt)
	return err
}

func UpdatePekerjaanAlumni(db *sql.DB, id int, pekerjaan *model.PekerjaanAlumni) error {
	query := `UPDATE pekerjaan_alumni SET alumni_id = $1, nama_perusahaan = $2, posisi_jabatan = $3, bidang_industri = $4, lokasi_kerja = $5, gaji_range = $6, tanggal_mulai_kerja = $7, tanggal_selesai_kerja = $8, status_pekerjaan = $9, deskripsi_pekerjaan = $10, updated_at = $11 WHERE id = $12`
	pekerjaan.UpdatedAt = time.Now()
	
	var tanggalSelesai interface{}
	if pekerjaan.TanggalSelesaiKerja != nil && !pekerjaan.TanggalSelesaiKerja.Time.IsZero() {
		tanggalSelesai = pekerjaan.TanggalSelesaiKerja.Time
	} else {
		tanggalSelesai = nil
	}
	
	_, err := db.Exec(query, 
		pekerjaan.AlumniID, 
		pekerjaan.NamaPerusahaan, 
		pekerjaan.PosisiJabatan, 
		pekerjaan.BidangIndustri, 
		pekerjaan.LokasiKerja, 
		pekerjaan.GajiRange, 
		pekerjaan.TanggalMulaiKerja.Time, 
		tanggalSelesai, 
		pekerjaan.StatusPekerjaan, 
		pekerjaan.DeskripsiPekerjaan, 
		pekerjaan.UpdatedAt, 
		id)
	return err
}

func DeletePekerjaanAlumni(db *sql.DB, id int) error {
	query := `DELETE FROM pekerjaan_alumni WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func SoftDeletePekerjaanAlumni(db *sql.DB, id int) error {
	query := `UPDATE pekerjaan_alumni SET is_delete = $1, updated_at = $2 WHERE id = $3`
	now := time.Now()
	_, err := db.Exec(query, now, now, id)
	return err
}

func GetPekerjaanAlumniByIDForSoftDelete(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	pekerjaan := new(model.PekerjaanAlumni)
	var tanggalMulai, tanggalSelesai sql.NullTime
	var isDelete sql.NullTime
	
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete FROM pekerjaan_alumni WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan, &pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja, &pekerjaan.GajiRange, &tanggalMulai, &tanggalSelesai, &pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt, &isDelete)
	if err != nil {
		return nil, err
	}
	
	if tanggalMulai.Valid {
		pekerjaan.TanggalMulaiKerja = model.CustomDate{Time: tanggalMulai.Time}
	}
	if tanggalSelesai.Valid {
		pekerjaan.TanggalSelesaiKerja = &model.CustomDate{Time: tanggalSelesai.Time}
	}
	if isDelete.Valid {
		pekerjaan.IsDelete = &isDelete.Time
	}
	
	return pekerjaan, nil
}

func GetPekerjaanAlumniWithPagination(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
		SELECT pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri, 
		       pa.lokasi_kerja, pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, 
		       pa.status_pekerjaan, pa.deskripsi_pekerjaan, pa.created_at, pa.updated_at, pa.is_delete
		FROM pekerjaan_alumni pa
		LEFT JOIN alumni a ON pa.alumni_id = a.id
		WHERE (pa.nama_perusahaan ILIKE $1 OR pa.posisi_jabatan ILIKE $1 OR pa.bidang_industri ILIKE $1 
		      OR pa.lokasi_kerja ILIKE $1 OR a.nama ILIKE $1) AND pa.is_delete IS NULL
		ORDER BY pa.%s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var tanggalMulai, tanggalSelesai sql.NullTime
		var isDelete sql.NullTime
		
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, 
			&p.LokasiKerja, &p.GajiRange, &tanggalMulai, &tanggalSelesai, 
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &isDelete)
		if err != nil {
			return nil, err
		}
		
		if tanggalMulai.Valid {
			p.TanggalMulaiKerja = model.CustomDate{Time: tanggalMulai.Time}
		}
		if tanggalSelesai.Valid {
			p.TanggalSelesaiKerja = &model.CustomDate{Time: tanggalSelesai.Time}
		}
		if isDelete.Valid {
			p.IsDelete = &isDelete.Time
		}
		
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func CountPekerjaanAlumni(db *sql.DB, search string) (int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pekerjaan_alumni pa
		LEFT JOIN alumni a ON pa.alumni_id = a.id
		WHERE (pa.nama_perusahaan ILIKE $1 OR pa.posisi_jabatan ILIKE $1 OR pa.bidang_industri ILIKE $1 
		      OR pa.lokasi_kerja ILIKE $1 OR a.nama ILIKE $1) AND pa.is_delete IS NULL
	`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}