package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoUploadRepository struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Path      string             `bson:"path"`
	Name      string             `bson:"name"`
	FileType  string             `bson:"file_type"`
	Type      string             `bson:"type"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
