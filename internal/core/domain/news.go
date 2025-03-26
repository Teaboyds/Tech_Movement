// entity //
package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title      string             `bson:"title" json:"title"`
	Detail     string             `bson:"detail" json:"detail"`
	Image      string             `bson:"image" json:"image"`
	CategoryID *Category          `bson:"category_id" json:"category_id"`
	Tag        []string           `bson:"tag" json:"tag"`
	CreatedAt  string             `bson:"created_at" json:"created_at"`
	UpdatedAt  string             `bson:"updated_at" json:"updated_at"`
}

type ErrResponse struct {
	Error string `json:"error"`
}
