package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"errors"
)

type MediaService struct {
	MediaRepo port.MediaRepository
}

func NewMediaService(MediaRepo port.MediaRepository) port.MediaService {
	return &MediaService{MediaRepo: MediaRepo}
}

func (med *MediaService) CreateMedia(media *domain.MediaRequest) error {

	if media.Category == "" {
		return errors.New("missing id in request body")
	}

	err := med.MediaRepo.CreateMedia(media)
	if err != nil {
		return err
	}

	return nil
}

func (med *MediaService) GetVideoHome() ([]*domain.VideoResponse, error) {
	video, err := med.MediaRepo.GetVideoHome()
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (med *MediaService) GetShortVideoHome() ([]*domain.ShortVideo, error) {
	short, err := med.MediaRepo.GetShortVideoHome()
	if err != nil {
		return nil, err
	}

	return short, nil
}
