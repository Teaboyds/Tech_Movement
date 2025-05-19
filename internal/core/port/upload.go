package port

import up "backend_tech_movement_hex/internal/core/domain"

type UploadRepository interface {
	SaveImage(file *up.UploadFile) error
	// GetFileByType(name string) []up.UploadFile
	GetFileByID(id string) (*up.UploadFile, error)
	GetAllFile() ([]up.UploadFile, error)
	ValidateImageIDs(ids []string) ([]string, error)
	GetFilesByIDs(ids []string) ([]up.UploadFile, error)
	DeleteFile(id string) error
	EnsureFileIndexs() error
	GetFilesByIDsVTest(ids []string) ([]*up.UploadFile, error)
}

type UploadService interface {
	UploadFile(file *up.UploadFileRequest) error
	// GetFileByType(name string) []up.UploadFile
	GetFileByID(id string) (*up.UploadFile, error)
	GetAllFile() ([]up.UploadFile, error)
	DeleteFile(id string) error
	GetFilesByIDsVTest(ids []string) ([]*up.UploadFile, error)
	ValidateImageIDs(ids []string) ([]string, error)
}
