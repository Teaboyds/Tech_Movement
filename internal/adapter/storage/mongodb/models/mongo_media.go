package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoMedia struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	URL        string             `bson:"url"`
	CategoryID primitive.ObjectID `bson:"category_id"`
	Status     bool               `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}
