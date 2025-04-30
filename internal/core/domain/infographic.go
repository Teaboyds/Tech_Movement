package domain

type Infographic struct {
	ID        string   `json:"id"`
	Image     string   `json:"image"`
	Title     string   `json:"title"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type InfographicRequest struct {
	Image    string   `json:"image"`
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Status   string   `json:"status"`
}

type InfographicRespose struct {
	ID        string             `json:"id"`
	Title     string             `json:"title"`
	Image     UploadFileResponse `json:"image"`
	CreatedAt string             `json:"created_at"`
}
