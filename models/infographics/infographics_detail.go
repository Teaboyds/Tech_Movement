package infographics

import "time"

type InfographicDetails struct {
	ID          uint      `gorm:"primaryKey"`
	DetailImage string    `gorm:"not null; type:text" json:"detail_image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	InfographicsID uint
}