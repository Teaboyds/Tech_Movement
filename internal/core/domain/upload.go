package domain

type UploadFile struct {
	ID        string `json:"id"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	FileType  string `json:"file_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UploadFileRequest struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	FileType string `form:"file_type" validate:"required,oneof=banner infographic news" json:"file_type"`
}

type UploadFileResponse struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	FileType string `json:"file_type"`
}

type UploadFileResponseHomePage struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Filetype string `json:"file_type"`
}
