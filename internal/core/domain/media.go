package domain

type Media struct {
	ID        string
	Title     string
	Content   string
	VideoURL  string
	Thumnail  string
	Category  string
	Tags      []string
	View      string
	Action    string
	CreatedAt string
	UpdatedAt string
}

type MediaRequest struct {
	Title    string   `json:"title" form:"title"`
	Content  string   `json:"content" form:"content"`
	VideoURL string   `json:"video_url" form:"video_url"`
	Thumnail string   `json:"thumnail_id" form:"thumnail_id"`
	Category string   `json:"category"  form:"category"`
	Tags     []string `json:"tags" form:"tags"`
	Action   string   `json:"action" form:"action"`
}

type VideoResponse struct {
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	URL       string           `json:"url"`
	Category  CategoryResponse `json:"category_id"`
	CreatedAt string           `json:"created_at"`
}

type MediaResponse struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	VideoUrl  string   `json:"video_url"`
	Thumnail  string   `json:"thumnail"`
	Category  Category `json:"category"`
	Tags      []string `json:"tags"`
	PageView  string   `json:"page_view"`
	Action    string   `json:"action"`
	CreatedAt string   `json:"created_at"`
}

type ShortVideo struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
