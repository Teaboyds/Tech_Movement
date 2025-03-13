package news

type News struct {
	ID             uint       `gorm:"primaryKey;"`
	Title          string     `json:"title" gorm:"not null;"`
	Image          string     `json:"image" gorm:"not null; type:text;"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
	NewsCategoryID uint       `gorm:"not null" json:"news_category_id"`
	NewsDetail     NewsDetail `gorm:"foreignKey:NewsID"`
	Tag            []*NewsTag `gorm:"many2many:news_newstag"`
}