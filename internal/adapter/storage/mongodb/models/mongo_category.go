package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoCategory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	CategoryType string             `bson:"category_type"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}
