package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type MediaService struct {
	MediaRepo       port.MediaRepository
	categoryRepo    port.CategoryRepository
	categoryService port.CategoryService
	fileService     port.UploadService
}

func NewMediaService(MediaRepo port.MediaRepository, categoryRepo port.CategoryRepository, categoryService port.CategoryService, fileService port.UploadService) port.MediaService {
	return &MediaService{MediaRepo: MediaRepo, categoryRepo: categoryRepo, categoryService: categoryService, fileService: fileService}
}

func (med *MediaService) CreateMedia(media *domain.Media) error {

	if media.CategoryID == "" {
		return errors.New("missing id in request body")
	}

	if media.View == "" {
		media.View = "0"
		_, err := strconv.Atoi(media.View)
		if err != nil {
			return fmt.Errorf("failed to parse View: %s", err)
		}
	}

	_, err := med.fileService.GetFileByID(media.ThumnailID)
	if err != nil {
		return fmt.Errorf("cannot fecth media thumnail : %v", err)
	}

	_, err = med.categoryRepo.GetByID(media.CategoryID)
	if err != nil {
		return fmt.Errorf("cannot fecth category thumnail :%v", err)
	}

	if media.Action == "" {
		media.Action = "On"
	} else {
		actionCheck := utils.IsValidAction(media.Action)
		if actionCheck == false {
			return fmt.Errorf("action must be 'on' or 'off'")
		}
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
			URL:      result.VideoURL,
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
			URL:   result.VideoURL,
		}

		responseShort = append(responseShort, resp)

	}

	return responseShort, nil
}
