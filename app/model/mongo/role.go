package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nama      string             `bson:"nama" json:"nama"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateRoleRequest struct {
	Nama string `bson:"nama" json:"nama" validate:"required"`
}

type UpdateRoleRequest struct {
	Nama string `bson:"nama" json:"nama" validate:"required"`
}

