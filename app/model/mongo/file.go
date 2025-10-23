package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AlumniID     primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	FileName     string             `bson:"file_name" json:"file_name"`
	OriginalName string             `bson:"original_name" json:"original_name"`
	FilePath     string             `bson:"file_path" json:"file_path"`
	FileSize     int64              `bson:"file_size" json:"file_size"`
	FileType     string             `bson:"file_type" json:"file_type"`
	Category     string             `bson:"category" json:"category"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type UploadFileRequest struct {
	AlumniID string `json:"alumni_id" validate:"required"`
}

type FileResponse struct {
	ID           string    `json:"id"`
	AlumniID     string    `json:"alumni_id"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
