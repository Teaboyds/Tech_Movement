package banners

import "time"

type BannersCategory struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Banners	[]Banners `gorm:"foreignKey:BannerCategoryID"`
}