package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoNews struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Title       string               `bson:"title"`
	Description string               `bson:"description"`
	Content     string               `bson:"content"`
	Image       []primitive.ObjectID `bson:"image"`
	CategoryID  primitive.ObjectID   `bson:"category_id"`
	Tag         []string             `bson:"tag"`
	Status      bool                 `bson:"status"`
	ContentType string               `bson:"content_type"`
	CreatedAt   time.Time            `bson:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at"`
}
