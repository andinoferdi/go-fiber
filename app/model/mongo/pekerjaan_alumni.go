package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlumniInfo struct {
	AlumniID    primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NIM         string             `bson:"nim" json:"nim"`
	Nama        string             `bson:"nama" json:"nama"`
	Email       string             `bson:"email" json:"email"`
}

type PekerjaanAlumni struct {
	ID                  primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	AlumniInfo          AlumniInfo          `bson:"alumni_info" json:"alumni_info"`
	NamaPerusahaan      string              `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string              `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string              `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string              `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           *string             `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
	TanggalMulaiKerja   time.Time           `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time          `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string              `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string             `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan,omitempty"`
	CreatedAt           time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time           `bson:"updated_at" json:"updated_at"`
}

type CreatePekerjaanAlumniRequest struct {
	AlumniInfo          AlumniInfo `bson:"alumni_info" json:"alumni_info" validate:"required"`
	NamaPerusahaan      string  `bson:"nama_perusahaan" json:"nama_perusahaan" validate:"required"`
	PosisiJabatan       string  `bson:"posisi_jabatan" json:"posisi_jabatan" validate:"required"`
	BidangIndustri      string  `bson:"bidang_industri" json:"bidang_industri" validate:"required"`
	LokasiKerja         string  `bson:"lokasi_kerja" json:"lokasi_kerja" validate:"required"`
	GajiRange           *string `bson:"gaji_range,omitempty" json:"gaji_range"`
	TanggalMulaiKerja   string  `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja *string `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string  `bson:"status_pekerjaan" json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanAlumniRequest struct {
	NamaPerusahaan      string  `bson:"nama_perusahaan" json:"nama_perusahaan" validate:"required"`
	PosisiJabatan       string  `bson:"posisi_jabatan" json:"posisi_jabatan" validate:"required"`
	BidangIndustri      string  `bson:"bidang_industri" json:"bidang_industri" validate:"required"`
	LokasiKerja         string  `bson:"lokasi_kerja" json:"lokasi_kerja" validate:"required"`
	GajiRange           *string `bson:"gaji_range,omitempty" json:"gaji_range"`
	TanggalMulaiKerja   string  `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja *string `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string  `bson:"status_pekerjaan" json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan"`
}

