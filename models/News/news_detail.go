package news

import "time"

type NewsDetail struct {
	ID        uint      `gorm:"primaryKey"`
	Detail    string    `gorm:"type:text" json:"detail"`
	NewsID    uint      `gorm:"unique; not null"`
	News      *News      `gorm:"foreignKey:NewsID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}