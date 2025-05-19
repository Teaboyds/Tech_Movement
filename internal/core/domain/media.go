package domain

type Media struct {
	ID         string
	Title      string
	Content    string
	URL        string
	CategoryID string
	Status     bool
	CreatedAt  string
	UpdatedAt  string
}

type MediaRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	URL      string `json:"url"`
	Category string `json:"category_id"`
	Status   bool   `json:"status"`
}

type VideoResponse struct {
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	URL       string           `json:"url"`
	Category  CategoryResponse `json:"category_id"`
	CreatedAt string           `json:"created_at"`
}

type ShortVideo struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
