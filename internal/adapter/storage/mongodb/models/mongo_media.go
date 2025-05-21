package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoMedia struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	VideoUrl  string             `bson:"video_url"`
	Thumnail  primitive.ObjectID `bson:"thumnail_id"`
	Category  primitive.ObjectID `bson:"category_id"`
	Tags      []string           `bson:"tags"`
	View      int                `bson:"view"`
	Action    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type CategoryMedia struct {
	ID   primitive.ObjectID `bson:"id"`
	Name string             `bson:"name"`
}
