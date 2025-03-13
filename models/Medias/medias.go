package medias

import (
	"time"
)

type Medias struct {
	ID uint `gorm:"primaryKey"`
	Title string `gorm:"not null" json:"title"`
	Detail string  `gorm:"not null; type:text" json:"detail"`
	Image string `gorm:"not null; type:text" json:"image"`
	URL string `gorm:"not null; type:text" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	MediaCategoryID uint
	MediaTags []*MediaTag `gorm:"many2many:medias_mediatag"`
}