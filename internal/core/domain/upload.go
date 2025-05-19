package domain

type UploadFile struct {
	ID        string `json:"id"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	FileType  string `json:"file_type"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UploadFileRequest struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	FileType string `form:"file_type" validate:"required,oneof=banner infographic news" json:"file_type"`
	Type     string `json:"type"`
}

type UploadFileResponse struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	FileType string `json:"file_type"`
	Type     string `json:"type"`
}

type UploadFileResponseHomePage struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	FileType string `json:"file_type"`
}
