package domain

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CategoryRequest struct {
	Name         string `json:"name"`
	CategoryType string `form:"category_type" validate:"required,oneof=news media" json:"category_type"`
}

type CategoryResponse struct {
	ID   string `json:"category_id"`
	Name string `json:"name"`
}
