package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoBanner struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Title       string               `bson:"name"`
	ContentType string               `bson:"content_type"`
	Status      bool                 `bson:"status"`
	CategoryID  primitive.ObjectID   `bson:"category_id"`
	ImageID     []primitive.ObjectID `bson:"image_id"`
	CreatedAt   time.Time            `bson:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at"`
}
