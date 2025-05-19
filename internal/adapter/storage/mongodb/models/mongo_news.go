package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoNews struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	ThumnailID  primitive.ObjectID   `bson:"thumnail_id"`
	Title       string               `bson:"title"`
	Description string               `bson:"description"`
	Content     string               `bson:"content"`
	ImageIDs    []primitive.ObjectID `bson:"image_ids"`
	CategoryID  primitive.ObjectID   `bson:"category_id"`
	Tags        []string             `bson:"tag"`
	Status      bool                 `bson:"status"`
	ContentType string               `bson:"content_type"`
	View        int                  `bson:"view"`
	CreatedAt   time.Time            `bson:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at"`
}
