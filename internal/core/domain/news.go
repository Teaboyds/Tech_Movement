// entity //
package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Detail        string             `bson:"detail" json:"detail"`
	Image         string             `bson:"image" json:"image"`
	CategoryID    *Category          `bson:"category_id" json:"category"`
	Tag           []string           `bson:"tag" json:"tag"`
	Status        bool               `bson:"status" json:"status"`
	ContentStatus string             `bson:"content_status" json:"content_status"` /* enum draft || publised || archived */
	CreatedAt     string             `bson:"created_at" json:"created_at"`
	UpdatedAt     string             `bson:"updated_at" json:"updated_at"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type NewsRequest struct {
	Title         string   `json:"title"`
	Detail        string   `json:"detail"`
	Image         string   `json:"image"`
	Category      string   `json:"category"`
	Tag           []string `json:"tag"`
	Status        bool     `json:"status"`
	ContentStatus string   `json:"content_status"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type UpdateNewsRequest struct {
	Title     string   `json:"title"`
	Detail    string   `json:"detail"`
	Image     string   `json:"image"`
	Category  string   `json:"category"`
	Tag       []string `json:"tag"`
	Status    bool     `json:"status"`
	UpdatedAt string   `json:"updated_at"`
}
