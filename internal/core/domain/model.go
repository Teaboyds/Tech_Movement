package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoModel struct {
	ID        primitive.ObjectID `json:"-" bson:"_id"`
	StrID     string             `json:"id" bson:"-"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
