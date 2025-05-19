package domain

type Media struct {
	ID         string
	Title      string
	Content    string
	VideoURL   string
	ThumnailID string
	CategoryID string
	Tags       []string
	View       string
	Action     string
	CreatedAt  string
	UpdatedAt  string
}

type MediaRequest struct {
	Title      string   `json:"title" form:"title"`
	Content    string   `json:"content" form:"content"`
	VideoURL   string   `json:"video_url" form:"video_url"`
	ThumnailID string   `json:"thumnail_id" form:"thumnail_id"`
	CategoryID string   `json:"category_id"  form:"category_id"`
	Tags       []string `json:"tags" form:"tags"`
	Action     string   `json:"action" form:"action"`
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
