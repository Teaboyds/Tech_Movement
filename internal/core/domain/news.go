// entity //
package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID         primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	Title      string                `bson:"title" json:"title"`
	Detail     string                `bson:"detail" json:"detail"`
	Image      ImageData             `bson:"image" json:"image"`
	CategoryID *Category             `bson:"category_id" json:"category"`
	Tag        *[]primitive.ObjectID `bson:"tag_id" json:"tag"`
	CreatedAt  string                `bson:"created_at" json:"created_at"`
	UpdatedAt  string                `bson:"updated_at" json:"updated_at"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type ImageData struct {
	ImagePath string `bson:"ImagePath" json:"ImagePath"`
	ImageName string `bson:"ImageName" json:"imageName"`
}
