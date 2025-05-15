package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	ul "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"
)

type UploadService struct {
	UploadRepo port.UploadRepository
}

func NewUploadService(UploadRepo port.UploadRepository) port.UploadService {
	return &UploadService{UploadRepo: UploadRepo}
}

func (ul *UploadService) UploadFile(file *ul.UploadFileRequest) error {

	input := &domain.UploadFile{
		Path:     file.Path,
		Name:     file.Name,
		FileType: file.FileType,
	}

	err := ul.UploadRepo.SaveImage(input)
	if err != nil {
		return err
	}

	return nil
}

func (ul *UploadService) GetFileByID(id string) (*ul.UploadFile, error) {

	file, err := ul.UploadRepo.GetFileByID(id)
	if err != nil {
		return nil, err
	}

	file.Path = utils.AttachBaseURLToImage(file.FileType, file.Path)

	return file, err
}

func (ul *UploadService) GetAllFile() ([]ul.UploadFile, error) {
	file, err := ul.UploadRepo.GetAllFile()
	if err != nil {
		return nil, err
	}

	for i, image := range file {
		file[i].Path = utils.AttachBaseURLToImage(image.FileType, image.Path)
	}

	return file, nil
}

func (ul *UploadService) DeleteFile(id string) error {

	file, err := ul.UploadRepo.GetFileByID(id)
	if err != nil {
		return fmt.Errorf("error by get by id")
	}

	fmt.Printf("file.ImageType: %v\n", file.FileType)

	if err := ul.UploadRepo.DeleteFile(id); err != nil {
		return err
	}

	fmt.Printf("file: %v\n", file)

	if file != nil {
		err := utils.DeleteFileInLocalStorage(file.FileType, file.Path)
		if err != nil {
			log.Println("Failed to delete image:", err)
		}
	}

	return nil
}

func (ul *UploadService) GetFilesByIDsVTest(ids []string) ([]ul.UploadFileResponse, error) {

	files, err := ul.UploadRepo.GetFilesByIDsVTest(ids)
	if err != nil {
		return nil, err
	}

	var resp []domain.UploadFileResponse
	for _, f := range files {
		resp = append(resp, domain.UploadFileResponse{
			Path:     f.Path,
			Name:     f.Name,
			FileType: f.FileType,
		})
	}

	return resp, nil
}

func (ul *UploadService) ValidateImageIDs(ids []string) ([]string, error) {
	return ul.UploadRepo.ValidateImageIDs(ids)
}
