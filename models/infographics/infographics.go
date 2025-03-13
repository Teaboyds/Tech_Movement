package infographics

import "time"

type Infographics struct {
	ID        string    `gorm:"primaryKey"`
	InfoImage string `gorm:"not null; type:text" json:"info_image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	InfographicDetails []InfographicDetails `gorm:"foreignKey:InfographicsID"`
}