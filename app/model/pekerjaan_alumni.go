package model

import (
	"encoding/json"
	"time"
)

type CustomDate struct {
	time.Time
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}
	
	if dateStr == "" || dateStr == "null" {
		cd.Time = time.Time{}
		return nil
	}
	
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	
	cd.Time = parsedTime
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(cd.Time.Format("2006-01-02"))
}

type PekerjaanAlumni struct {
	ID                  int         `json:"id" db:"id"`
	AlumniID           int         `json:"alumni_id" db:"alumni_id"`
	NamaPerusahaan     string      `json:"nama_perusahaan" db:"nama_perusahaan"`
	PosisiJabatan      string      `json:"posisi_jabatan" db:"posisi_jabatan"`
	BidangIndustri     string      `json:"bidang_industri" db:"bidang_industri"`
	LokasiKerja        string      `json:"lokasi_kerja" db:"lokasi_kerja"`
	GajiRange          *string     `json:"gaji_range" db:"gaji_range"`
	TanggalMulaiKerja  CustomDate  `json:"tanggal_mulai_kerja" db:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *CustomDate `json:"tanggal_selesai_kerja" db:"tanggal_selesai_kerja"`
	StatusPekerjaan    string      `json:"status_pekerjaan" db:"status_pekerjaan"`
	DeskripsiPekerjaan *string     `json:"deskripsi_pekerjaan" db:"deskripsi_pekerjaan"`
	CreatedAt          time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at" db:"updated_at"`
	IsDelete           *time.Time  `json:"is_delete" db:"is_delete"`
}