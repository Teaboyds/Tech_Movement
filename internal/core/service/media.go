package service

import (
	dm "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
)

type MediaService struct {
	MediaRepo port.MediaRepository
}

func NewMediaService(MediaRepo port.MediaRepository) port.MediaService {
	return &MediaService{MediaRepo: MediaRepo}
}

func (m *MediaService) CreateMedia(media *dm.Media) error {
	return m.MediaRepo.SaveMedia(media)
}

func (m *MediaService) GetLastMedia() ([]dm.Media, error) {

	media, err := m.MediaRepo.GetLastMedia()
	if err != nil {
		return nil, err
	}

	for i := range media {
		utils.AttachBaseURLToMediaImg(&media[i])
	}

	return media, nil
}
