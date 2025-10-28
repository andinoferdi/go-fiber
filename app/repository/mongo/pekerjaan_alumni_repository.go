package repository

import (
	"context"
	"errors"
	model "go-fiber/app/model/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IPekerjaanAlumniRepository interface {
	CreatePekerjaanAlumni(ctx context.Context, pekerjaan *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error)
	FindPekerjaanAlumniByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error)
	FindAllPekerjaanAlumni(ctx context.Context) ([]model.PekerjaanAlumni, error)
	FindPekerjaanAlumniByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error)
	UpdatePekerjaanAlumni(ctx context.Context, id string, pekerjaan *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error)
	DeletePekerjaanAlumni(ctx context.Context, id string) error
}

type PekerjaanAlumniRepository struct {
	collection *mongo.Collection
}

func NewPekerjaanAlumniRepository(db *mongo.Database) IPekerjaanAlumniRepository {
	return &PekerjaanAlumniRepository{
		collection: db.Collection("pekerjaan_alumni"),
	}
}

func (r *PekerjaanAlumniRepository) CreatePekerjaanAlumni(ctx context.Context, pekerjaan *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	pekerjaan.ID = primitive.NilObjectID
	pekerjaan.CreatedAt = time.Now()
	pekerjaan.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, pekerjaan)
	if err != nil {
		return nil, err
	}

	pekerjaan.ID = result.InsertedID.(primitive.ObjectID)
	return pekerjaan, nil
}

func (r *PekerjaanAlumniRepository) FindPekerjaanAlumniByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	var pekerjaan model.PekerjaanAlumni
	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&pekerjaan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &pekerjaan, nil
}

func (r *PekerjaanAlumniRepository) FindAllPekerjaanAlumni(ctx context.Context) ([]model.PekerjaanAlumni, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaanList []model.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaanList); err != nil {
		return nil, err
	}

	return pekerjaanList, nil
}

func (r *PekerjaanAlumniRepository) FindPekerjaanAlumniByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error) {
	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return nil, errors.New("alumni ID tidak valid")
	}

	filter := bson.M{"alumni_info.alumni_id": objID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaanList []model.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaanList); err != nil {
		return nil, err
	}

	return pekerjaanList, nil
}

func (r *PekerjaanAlumniRepository) UpdatePekerjaanAlumni(ctx context.Context, id string, pekerjaan *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"alumni_info":          pekerjaan.AlumniInfo,
			"nama_perusahaan":      pekerjaan.NamaPerusahaan,
			"posisi_jabatan":       pekerjaan.PosisiJabatan,
			"bidang_industri":      pekerjaan.BidangIndustri,
			"lokasi_kerja":         pekerjaan.LokasiKerja,
			"gaji_range":           pekerjaan.GajiRange,
			"tanggal_mulai_kerja":  pekerjaan.TanggalMulaiKerja,
			"tanggal_selesai_kerja": pekerjaan.TanggalSelesaiKerja,
			"status_pekerjaan":     pekerjaan.StatusPekerjaan,
			"deskripsi_pekerjaan":  pekerjaan.DeskripsiPekerjaan,
			"updated_at":           pekerjaan.UpdatedAt,
		},
	}

	result := r.collection.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}

	return r.FindPekerjaanAlumniByID(ctx, id)
}

func (r *PekerjaanAlumniRepository) DeletePekerjaanAlumni(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID tidak valid")
	}

	filter := bson.M{"_id": objID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("data pekerjaan alumni tidak ditemukan")
	}

	return nil
}

