package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"fmt"
)

type InfographicService struct {
	InfoRepo    port.InfographicRepository
	fileService port.UploadService
	cateRepo    port.CategoryService
}

func NewInfographicService(InfoRepo port.InfographicRepository, fileService port.UploadService, cateRepo port.CategoryService) port.InfographicService {
	return &InfographicService{InfoRepo: InfoRepo, fileService: fileService, cateRepo: cateRepo}
}

func (ip *InfographicService) CreateInfo(info *domain.Infographic) error {

	_, err := ip.cateRepo.GetByID(info.Category)
	if err != nil {
		return err
	}

	_, err = ip.fileService.GetFileByID(info.Image)
	if err != nil {
		return err
	}

	if info.PageView == "" {
		info.PageView = "0"
	}

	err = ip.InfoRepo.CreateInfo(info)
	if err != nil {
		return err
	}

	return nil
}

func (ip *InfographicService) GetInfographic(id string) (*domain.Infographic, error) {
	resp, err := ip.InfoRepo.Retrive(id)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (ip *InfographicService) GetInfographics(cateId, sort, view, limit, page string) ([]*domain.Infographic, error) {

	if limit == "" {
		limit = "10"
	}

	if page == "" {
		page = "1"
	}

	infographic, err := ip.InfoRepo.RetrivesInfographic(cateId, sort, view, limit, page)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch media in media service")
	}

	return infographic, err
}

// func (ip *InfographicService) GetInfoHome() ([]*domain.InfographicRespose, error) {

// 	info, err := ip.InfoRepo.GetInfoHome()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var response []*domain.InfographicRespose
// 	for _, infoo := range info {

// 		fmt.Printf("infoo.Image: %v\n", infoo.Image)

// 		uploadFile, err := ip.fileService.GetFileByID(infoo.Image)
// 		if err != nil {
// 			log.Println(err)
// 			return nil, fmt.Errorf("file not found")
// 		}

// 		image := &domain.UploadFileResponse{
// 			ID:       uploadFile.ID,
// 			Path:     uploadFile.Path,
// 			Name:     uploadFile.Name,
// 			FileType: uploadFile.FileType,
// 			Type:     uploadFile.Type,
// 		}

// 		fmt.Printf("image.FileType: %v\n", image.FileType)
// 		fmt.Printf("image.Path: %v\n", image.Path)

// 		resp := &domain.InfographicRespose{
// 			ID:        infoo.ID,
// 			Title:     infoo.Title,
// 			Image:     *image,
// 			CreatedAt: infoo.CreatedAt,
// 		}

// 		response = append(response, resp)

// 	}

// 	return response, nil
// }
