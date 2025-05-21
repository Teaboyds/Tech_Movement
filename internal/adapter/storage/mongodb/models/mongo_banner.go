package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoBanner struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	DesktopImage ImageInfo          `bson:"desktop_image"`
	MobileImage  ImageInfo          `bson:"mobile_image"`
	Status       StatusType         `bson:"status"`
	LinkUrl      string             `bson:"link_url"`
	Action       string             `bson:"action"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type MongoBannerV2 struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	DesktopImage string             `bson:"desktop_image"`
	MobileImage  string             `bson:"mobile_image"`
	Status       StatusType         `bson:"status"`
	LinkUrl      string             `bson:"link_url"`
	Action       string             `bson:"action"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type ImageInfo struct {
	Path     string `bson:"path"`
	Name     string `bson:"name"`
	FileType string `bson:"file_type"`
	Type     string `bson:"type"`
}

type StatusType struct {
	Home        bool `bson:"home"`
	Media       bool `bson:"media"`
	News        bool `bson:"news"`
	Infographic bool `bson:"infographic"`
}
