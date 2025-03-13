package medias

import (
	"time"
)

type MediaCategory struct {
	ID          uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null; default:'Uncategorized'" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Medias []Medias`gorm:"foreignKey:MediaCategoryID"`
}