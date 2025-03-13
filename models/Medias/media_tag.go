package medias

import (
	"time"
)

type MediaTag struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}