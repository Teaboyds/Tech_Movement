package domain

type Media struct {
	ID         string `json:"_id,omitempty"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	URL        string `json:"url"`
	CategoryID string `json:"category_id"`
	Status     bool   `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
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
