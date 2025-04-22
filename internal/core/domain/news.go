// entity //
package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type News struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Detail        string             `bson:"detail" json:"detail"`
	Image         string             `bson:"image" json:"image"`
	CategoryID    *Category          `bson:"category_id" json:"category"`
	Tag           []string           `bson:"tag" json:"tag"`
	Status        bool               `bson:"status" json:"status"`
	ContentStatus string             `bson:"content_status" json:"content_status"` /* enum draft || publised || archived */
	ContentType   string             `bson:"content_type" json:"content_type"`     /* enum general || breaking || video */
	CreatedAt     string             `bson:"created_at" json:"created_at"`
	UpdatedAt     string             `bson:"updated_at" json:"updated_at"`
}

////////////////////////////////// Response // Request // Models /////////////////////////////////////////////////

type NewsRequest struct {
	Title         string `json:"title" validate:"required"`
	Detail        string `json:"detail" validate:"required"`
	Image         string `json:"image"`
	Category      string `json:"category"`
	Tag           string `json:"tag" validate:"required,min=1,required"`
	Status        string `json:"status" validate:"required"`
	ContentStatus string `form:"content_status" validate:"required,oneof=draft published archived"`
	ContentType   string `form:"content_type" validate:"required,oneof=general breaking video"`
}

type UpdateNewsRequestResponse struct {
	Title         string `json:"title"`
	Detail        string `json:"detail"`
	Image         string `json:"image"`
	Category      string `json:"category"`
	Tag           string `json:"tag"`
	Status        string `json:"status"`
	ContentStatus string `form:"content_status"`
	ContentType   string `form:"content_type"`
	UpdatedAt     string `json:"updated_at"`
}

// News Home::LastedNews//
type HomePageLastedNewResponse struct {
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Image     string `json:"image"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
}

// News Home::Technology//
type NewsHomePageResponse struct {
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Image     string `json:"image"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

////////////////////////////////// Response Models /////////////////////////////////////////////////
