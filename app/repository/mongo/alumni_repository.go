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

type IAlumniRepository interface {
	CreateAlumni(ctx context.Context, alumni *model.Alumni) (*model.Alumni, error)
	FindAlumniByID(ctx context.Context, id string) (*model.Alumni, error)
	FindAlumniByEmail(ctx context.Context, email string) (*model.Alumni, error)
	FindAlumniByNIM(ctx context.Context, nim string) (*model.Alumni, error)
	FindAllAlumni(ctx context.Context) ([]model.Alumni, error)
	UpdateAlumni(ctx context.Context, id string, alumni *model.Alumni) (*model.Alumni, error)
	DeleteAlumni(ctx context.Context, id string) error
}

type AlumniRepository struct {
	collection     *mongo.Collection
	roleCollection *mongo.Collection
}

func NewAlumniRepository(db *mongo.Database) IAlumniRepository {
	return &AlumniRepository{
		collection:     db.Collection("alumni"),
		roleCollection: db.Collection("roles"),
	}
}

func (r *AlumniRepository) populateRole(ctx context.Context, alumni *model.Alumni) error {
	if alumni.RoleID.IsZero() {
		return nil
	}

	var role model.Role
	filter := bson.M{"_id": alumni.RoleID}
	err := r.roleCollection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}

	alumni.Role = &role
	return nil
}

func (r *AlumniRepository) CreateAlumni(ctx context.Context, alumni *model.Alumni) (*model.Alumni, error) {
	alumni.ID = primitive.NilObjectID
	alumni.CreatedAt = time.Now()
	alumni.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, alumni)
	if err != nil {
		return nil, err
	}

	alumni.ID = result.InsertedID.(primitive.ObjectID)

	if err := r.populateRole(ctx, alumni); err != nil {
		return alumni, nil
	}

	return alumni, nil
}

func (r *AlumniRepository) FindAlumniByID(ctx context.Context, id string) (*model.Alumni, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	var alumni model.Alumni
	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&alumni)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	if err := r.populateRole(ctx, &alumni); err != nil {
		return &alumni, nil
	}

	return &alumni, nil
}

func (r *AlumniRepository) FindAlumniByEmail(ctx context.Context, email string) (*model.Alumni, error) {
	var alumni model.Alumni
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&alumni)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	if err := r.populateRole(ctx, &alumni); err != nil {
		return nil, err
	}

	return &alumni, nil
}

func (r *AlumniRepository) FindAlumniByNIM(ctx context.Context, nim string) (*model.Alumni, error) {
	var alumni model.Alumni
	filter := bson.M{"nim": nim}
	err := r.collection.FindOne(ctx, filter).Decode(&alumni)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	if err := r.populateRole(ctx, &alumni); err != nil {
		return &alumni, nil
	}

	return &alumni, nil
}

func (r *AlumniRepository) FindAllAlumni(ctx context.Context) ([]model.Alumni, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumniList []model.Alumni
	if err = cursor.All(ctx, &alumniList); err != nil {
		return nil, err
	}

	for i := range alumniList {
		if err := r.populateRole(ctx, &alumniList[i]); err != nil {
			continue
		}
	}

	return alumniList, nil
}

func (r *AlumniRepository) UpdateAlumni(ctx context.Context, id string, alumni *model.Alumni) (*model.Alumni, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"nama":        alumni.Nama,
			"jurusan":     alumni.Jurusan,
			"angkatan":    alumni.Angkatan,
			"tahun_lulus": alumni.TahunLulus,
			"email":       alumni.Email,
			"no_telepon":  alumni.NoTelepon,
			"alamat":      alumni.Alamat,
			"role_id":     alumni.RoleID,
			"updated_at":  time.Now(),
		},
	}

	result := r.collection.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}

	return r.FindAlumniByID(ctx, id)
}

func (r *AlumniRepository) DeleteAlumni(ctx context.Context, id string) error {
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
		return errors.New("Data alumni tidak ditemukan")
	}

	return nil
}

