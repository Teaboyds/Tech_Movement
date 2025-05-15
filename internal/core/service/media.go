package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"errors"
	"fmt"
	"log"
)

type MediaService struct {
	MediaRepo       port.MediaRepository
	categoryRepo    port.CategoryRepository
	categoryService port.CategoryService
}

func NewMediaService(MediaRepo port.MediaRepository, categoryRepo port.CategoryRepository, categoryService port.CategoryService) port.MediaService {
	return &MediaService{MediaRepo: MediaRepo, categoryRepo: categoryRepo, categoryService: categoryService}
}

func (med *MediaService) CreateMedia(media *domain.MediaRequest) error {

	if media.Category == "" {
		return errors.New("missing id in request body")
	}

	_, err := med.categoryRepo.GetByID(media.Category)
	if err != nil {
		return err
	}

	inputMedia := &domain.Media{
		Title:      media.Title,
		Content:    media.Content,
		URL:        media.URL,
		CategoryID: media.Category,
		Status:     media.Status,
	}

	err = med.MediaRepo.CreateMedia(inputMedia)
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

	if category == nil {
		return nil, fmt.Errorf("category 'short vdo' not found")
	}

	video, err := med.MediaRepo.GetVideoHome(category.ID)
	if err != nil {
		return nil, err
	}

	categoryIDMap := make(map[string]struct{})

	for _, vdo := range video {
		if vdo.CategoryID != "" {
			categoryIDMap[vdo.CategoryID] = struct{}{}
		}
	}

	categoryIDs := keysFromMap(categoryIDMap)

	categories, err := med.categoryService.GetByIDs(categoryIDs)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]domain.CategoryResponse)
	for _, ca := range categories {
		categoryMap[ca.ID] = *ca
	}

	var responseVDO []*domain.VideoResponse

	for _, result := range video {

		var categoryResponse domain.CategoryResponse
		if result.CategoryID != "" {
			if cat, ok := categoryMap[result.CategoryID]; ok {
				categoryResponse = cat
			}
		}

		resp := &domain.VideoResponse{
			Title:    result.Title,
			Content:  result.Content,
			URL:      result.URL,
			Category: categoryResponse,
		}

		responseVDO = append(responseVDO, resp)

	}

	return responseVDO, nil
}

func (med *MediaService) GetShortVideoHome() ([]*domain.ShortVideo, error) {

	categories, err := med.categoryRepo.GetByName("Short Video")
	if err != nil {
		return nil, fmt.Errorf("error fetching category: %w", err)
	}

	short, err := med.MediaRepo.GetShortVideoHome(categories.ID)
	if err != nil {
		return nil, err
	}

	var responseShort []*domain.ShortVideo

	for _, result := range short {

		resp := &domain.ShortVideo{
			Title: result.Title,
			URL:   result.URL,
		}

		responseShort = append(responseShort, resp)

	}

	return responseShort, nil
}
