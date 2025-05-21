package domain

type Infographic struct {
	ID        string
	Image     string
	Title     string
	Category  string
	Tags      []string
	Status    string
	PageView  string
	CreatedAt string
	UpdatedAt string
}

// DTO //
type InfographicRequestDTO struct {
	Image    string   `json:"image"`
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Status   string   `json:"status"`
}

type InfographicRespose struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Category  Category   `json:"category"`
	Image     UploadFile `json:"image"`
	Tags      []string   `json:"tags"`
	Status    string     `json:"status"`
	PageView  string     `json:"page_view"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}
