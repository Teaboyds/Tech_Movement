package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Abstract  string             `bson:"abstract" json:"abstract"`
	Image     string             `bson:"image" json:"image"`
	Url       string             `bson:"url" json:"url"`
	Category  *Category          `bson:"category_id" json:"category_id"`
	Tag       []string           `bson:"tag" json:"tag"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type MediaRequest struct {
	Title    string ` json:"title"`
	Abstract string ` json:"abstract"`
	Image    string ` json:"image"`
	Url      string ` json:"url"`
	Category string ` json:"category"`
	Tag      string ` json:"tag"`
}

type HomeMediaResponse struct {
	Title     string    `json:"title"`
	Abstract  string    ` json:"abstract"`
	Image     string    ` json:"image"`
	Url       string    ` json:"url"`
	Category  string    ` json:"category"`
	CreatedAt time.Time `json:"created_at"`
}
