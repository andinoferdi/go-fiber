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

type IRoleRepository interface {
	CreateRole(ctx context.Context, role *model.Role) (*model.Role, error)
	FindRoleByID(ctx context.Context, id string) (*model.Role, error)
	FindRoleByName(ctx context.Context, nama string) (*model.Role, error)
	FindAllRoles(ctx context.Context) ([]model.Role, error)
	UpdateRole(ctx context.Context, id string, role *model.Role) (*model.Role, error)
	DeleteRole(ctx context.Context, id string) error
}

type RoleRepository struct {
	collection *mongo.Collection
}

func NewRoleRepository(db *mongo.Database) IRoleRepository {
	return &RoleRepository{
		collection: db.Collection("roles"),
	}
}

func (r *RoleRepository) CreateRole(ctx context.Context, role *model.Role) (*model.Role, error) {
	role.ID = primitive.NilObjectID
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		return nil, err
	}

	role.ID = result.InsertedID.(primitive.ObjectID)
	return role, nil
}

func (r *RoleRepository) FindRoleByID(ctx context.Context, id string) (*model.Role, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	var role model.Role
	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) FindRoleByName(ctx context.Context, nama string) (*model.Role, error) {
	var role model.Role
	filter := bson.M{"nama": nama}
	err := r.collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) FindAllRoles(ctx context.Context) ([]model.Role, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []model.Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *RoleRepository) UpdateRole(ctx context.Context, id string, role *model.Role) (*model.Role, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"nama":       role.Nama,
			"updated_at": time.Now(),
		},
	}

	result := r.collection.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}

	return r.FindRoleByID(ctx, id)
}

func (r *RoleRepository) DeleteRole(ctx context.Context, id string) error {
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
		return errors.New("Data role tidak ditemukan")
	}

	return nil
}

