package news

type NewsCategory struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"default:'Uncategorized'; not null" json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	News []News `gorm:"foreignKey:NewsCategoryID"`
}