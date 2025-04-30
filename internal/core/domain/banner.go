package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Banner struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}
