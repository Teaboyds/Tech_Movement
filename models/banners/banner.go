package banners

import (
	"time"
)

type Banners struct {
	ID uint `gorm:"primaryKey"`
	Image string `gorm:"type:text; not null" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	BannerCategoryID uint `gorm:"not null"`
}
