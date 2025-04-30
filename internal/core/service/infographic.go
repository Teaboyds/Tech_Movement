package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
)

type InfographicService struct {
	InfoRepo port.InfographicRepository
}

func NewInfographicService(InfoRepo port.InfographicRepository) port.InfographicService {
	return &InfographicService{InfoRepo: InfoRepo}
}

func (ip *InfographicService) CreateInfo(info *domain.InfographicRequest) error {

	err := ip.InfoRepo.CreateInfo(info)
	if err != nil {
		return err
	}

	return nil
}

func (ip *InfographicService) GetInfoHome() ([]domain.InfographicRespose, error) {
	info, err := ip.InfoRepo.GetInfoHome()
	if err != nil {
		return nil, err
	}

	for i := range info {
		info[i].Image.Path = utils.AttachBaseURLToImage(info[i].Image.FileType, info[i].Image.Path)

	}

	return info, nil
}
