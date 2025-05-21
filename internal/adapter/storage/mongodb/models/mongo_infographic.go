package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoInfographic struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Image     primitive.ObjectID `bson:"image"`
	Title     string             `bson:"title"`
	Category  primitive.ObjectID `bson:"category"`
	Tags      []string           `bson:"tags"`
	Status    bool               `bson:"status"`
	PageView  int                `json:"page_view"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
