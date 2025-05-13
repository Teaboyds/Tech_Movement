package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"errors"
	"log"
)

type MediaService struct {
	MediaRepo    port.MediaRepository
	categoryRepo port.CategoryRepository
}

func NewMediaService(MediaRepo port.MediaRepository, categoryRepo port.CategoryRepository) port.MediaService {
	return &MediaService{MediaRepo: MediaRepo, categoryRepo: categoryRepo}
}

func (med *MediaService) CreateMedia(media *domain.MediaRequest) error {

	if media.Category == "" {
		return errors.New("missing id in request body")
	}

	_, err := med.categoryRepo.GetByID(media.Category)
	if err != nil {
		return err
	}

	err = med.MediaRepo.CreateMedia(media)
	if err != nil {
		return err
	}

	return nil
}

func (med *MediaService) GetVideoHome() ([]*domain.VideoResponse, error) {

	category, err := med.categoryRepo.GetByName("Short Video")
	if err != nil {
		log.Println("error fetching short_video category:", err)
		return nil, err
	}

	video, err := med.MediaRepo.GetVideoHome(category.ID)
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
