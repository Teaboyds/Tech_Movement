package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoMedia struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	VideoUrl   string             `bson:"video_url"`
	ThumnailID primitive.ObjectID `bson:"thumnail_id"`
	CategoryID primitive.ObjectID `bson:"category_id"`
	Tags       []string           `bson:"tags"`
	View       int                `bson:"view"`
	Action     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}
