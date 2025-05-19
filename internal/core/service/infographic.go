package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"fmt"
	"log"
	"strconv"
)

type InfographicService struct {
	InfoRepo    port.InfographicRepository
	fileService port.UploadService
	cateRepo    port.CategoryService
}

func NewInfographicService(InfoRepo port.InfographicRepository, fileService port.UploadService, cateRepo port.CategoryService) port.InfographicService {
	return &InfographicService{InfoRepo: InfoRepo, fileService: fileService, cateRepo: cateRepo}
}

func (ip *InfographicService) CreateInfo(info *domain.InfographicRequest) error {

	statusBool, err := strconv.ParseBool(info.Status)
	if err != nil {
		return fmt.Errorf("cannot phase bool")
	}

	_, err = ip.cateRepo.GetByID(info.Category)
	if err != nil {
		return err
	}

	_, err = ip.fileService.GetFileByID(info.Image)
	if err != nil {
		return err
	}

	input := &domain.Infographic{
		Image:    info.Image,
		Title:    info.Title,
		Category: info.Category,
		Tags:     info.Tags,
		Status:   statusBool,
	}

	err = ip.InfoRepo.CreateInfo(input)
	if err != nil {
		return err
	}

	return nil
}

func (ip *InfographicService) GetInfoHome() ([]*domain.InfographicRespose, error) {

	info, err := ip.InfoRepo.GetInfoHome()
	if err != nil {
		return nil, err
	}

	var response []*domain.InfographicRespose
	for _, infoo := range info {

		fmt.Printf("infoo.Image: %v\n", infoo.Image)

		uploadFile, err := ip.fileService.GetFileByID(infoo.Image)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("file not found")
		}

		image := &domain.UploadFileResponse{
			ID:       uploadFile.ID,
			Path:     uploadFile.Path,
			Name:     uploadFile.Name,
			FileType: uploadFile.FileType,
			Type:     uploadFile.Type,
		}

		fmt.Printf("image.FileType: %v\n", image.FileType)
		fmt.Printf("image.Path: %v\n", image.Path)

		resp := &domain.InfographicRespose{
			ID:        infoo.ID,
			Title:     infoo.Title,
			Image:     *image,
			CreatedAt: infoo.CreatedAt,
		}

		response = append(response, resp)

	}

	return response, nil
}
