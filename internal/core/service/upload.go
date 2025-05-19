package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	ul "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"
	"strings"
)

type UploadService struct {
	UploadRepo port.UploadRepository
}

func NewUploadService(UploadRepo port.UploadRepository) port.UploadService {
	return &UploadService{UploadRepo: UploadRepo}
}

// upload ภาพ //
func (ul *UploadService) UploadFile(file *ul.UploadFileRequest) error {

	namskun := strings.TrimPrefix(file.Type, ".")
	alt := strings.ReplaceAll(file.Name, " ", "_")

	input := &domain.UploadFile{
		Path:     file.Path,
		Name:     alt,
		FileType: file.FileType,
		Type:     namskun,
	}

	err := ul.UploadRepo.SaveImage(input)
	if err != nil {
		return err
	}

	return nil
}

// ใช้ Get ตอนต้องการดึงภาพ By ID //
func (ul *UploadService) GetFileByID(id string) (*ul.UploadFile, error) {

	file, err := ul.UploadRepo.GetFileByID(id)
	if err != nil {
		return nil, err
	}

	file.Path = utils.AttachBaseURLToImage(file.FileType, file.Path)

	return file, err
}

// ใช้ดึง file ภาพทั้งหมด //
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

// ลบภาพ //
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

// ใช้ Get ตอนต้องการดึงข่าวหลายตัว แล้วก็เอาไว้ run id จาก array ที่รับเข้ามาแล้ว returm ส่งออกไปเป็น array //
func (ul *UploadService) GetFilesByIDsVTest(ids []string) ([]*ul.UploadFile, error) {

	files, err := ul.UploadRepo.GetFilesByIDsVTest(ids)
	if err != nil {
		return nil, err
	}

	var resp []*domain.UploadFile
	for _, f := range files {
		resp = append(resp, &domain.UploadFile{
			ID:       f.ID,
			Name:     f.Name,
			Path:     f.Path,
			FileType: f.FileType,
		})
	}

	for _, item := range resp {
		item.Path = utils.AttachBaseURLToImage(item.FileType, item.Path)
	}

	return resp, nil
}

// ใช้เช็คว่า []ids มีอยู่ใน database ไหมแล้วคืนไปเป็น []string เช่รเดิม  **อาจจะลบ //
func (ul *UploadService) ValidateImageIDs(ids []string) ([]string, error) {
	return ul.UploadRepo.ValidateImageIDs(ids)
}
