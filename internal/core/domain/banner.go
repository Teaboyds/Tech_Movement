package domain

// domain หลัก //
type Banner struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	ContentType string   `json:"content_type"`
	Status      bool     `json:"status"`
	Category    string   `json:"category"`
	Img         []string `json:"image"`
}

// รับ request //
type BannerReq struct {
	Title       string   `json:"title"`
	ContentType string   `json:"content_type" form:"content_type"`
	Status      bool     `json:"status"`
	Category    string   `json:"category"`
	Img         []string `json:"image"`
}

// Response //
type BannerClient struct {
	Title       string            `json:"title"`
	ContentType string            `json:"content_type"`
	Status      bool              `json:"status"`
	Category    *CategoryResponse `json:"category"`
	Images      []UploadFileResponse
}
