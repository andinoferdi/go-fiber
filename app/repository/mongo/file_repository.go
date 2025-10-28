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

type IFileRepository interface {
	CreateFile(ctx context.Context, file *model.File) (*model.File, error)
	FindFilesByAlumniID(ctx context.Context, alumniID string) ([]model.File, error)
	FindFileByID(ctx context.Context, id string) (*model.File, error)
	DeleteFile(ctx context.Context, id string) error
}

type FileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) IFileRepository {
	return &FileRepository{
		collection: db.Collection("files"),
	}
}

func (r *FileRepository) CreateFile(ctx context.Context, file *model.File) (*model.File, error) {
	file.ID = primitive.NilObjectID
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		return nil, err
	}

	file.ID = result.InsertedID.(primitive.ObjectID)
	return file, nil
}

func (r *FileRepository) FindFilesByAlumniID(ctx context.Context, alumniID string) ([]model.File, error) {
	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return nil, errors.New("alumni ID tidak valid")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"alumni_info.alumni_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []model.File
	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *FileRepository) FindFileByID(ctx context.Context, id string) (*model.File, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	var file model.File
	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) DeleteFile(ctx context.Context, id string) error {
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
		return errors.New("file tidak ditemukan")
	}

	return nil
}
