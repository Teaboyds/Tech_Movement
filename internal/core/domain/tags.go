package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tags struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type TagsResponse struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}
