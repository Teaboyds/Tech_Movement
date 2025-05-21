package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
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

	if media.Category == "" {
		return fmt.Errorf("please input your category id")
	}

	if media.View == "" {
		media.View = "0"
		_, err := strconv.Atoi(media.View)
		if err != nil {
			return fmt.Errorf("failed to parse View: %s", err)
		}
	}

	_, err := med.fileService.GetFileByID(media.Thumnail)
	if err != nil {
		return fmt.Errorf("cannot fecth media thumnail : %v", err)
	}

	_, err = med.categoryRepo.GetByID(media.Category)
	if err != nil {
		return fmt.Errorf("cannot fecth category thumnail :%v", err)
	}

	if media.Action == "" {
		media.Action = "On"
	} else {
		actionCheck := utils.IsValidAction(media.Action)
		if !actionCheck {
			return fmt.Errorf("action must be 'on' or 'off'")
		}
	}

	err = med.MediaRepo.CreateMedia(media)
	if err != nil {
		return err
	}

	return nil
}

func (med *MediaService) GetMedias(cateId, sort, view, limit, page string) ([]*domain.Media, error) {

	if cateId == "" {
		return nil, fmt.Errorf("please input category id")
	}

	if sort == "" {
	} else if sort != "newest" && sort != "oldest" {
		return nil, fmt.Errorf("sort must be 'newest' or 'oldest'")
	}

	if view == "" {
	} else if view != "asc" && view != "desc" {
		return nil, fmt.Errorf("view must be 'asc' or 'desc'")
	}

	if limit == "" {
		limit = "10"
	}

	if page == "" {
		page = "1"
	}

	media, err := med.MediaRepo.RetrivesMedia(cateId, sort, view, limit, page)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch media in media service")
	}

	return media, err
}

func (med *MediaService) GetMedia(id string) (*domain.Media, error) {
	resp, err := med.MediaRepo.RetriveMedia(id)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// func (med *MediaService) GetVideoHome() ([]*domain.VideoResponse, error) {

// 	category, err := med.categoryRepo.GetByName("Short Video")
// 	if err != nil {
// 		log.Println("error fetching short_video category:", err)
// 		return nil, err
// 	}

// 	if category == nil {
// 		return nil, fmt.Errorf("category 'short vdo' not found")
// 	}

// 	video, err := med.MediaRepo.GetVideoHome(category.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	categoryIDMap := make(map[string]struct{})

// 	for _, vdo := range video {
// 		if vdo.CategoryID != "" {
// 			categoryIDMap[vdo.CategoryID] = struct{}{}
// 		}
// 	}

// 	categoryIDs := keysFromMap(categoryIDMap)

// 	categories, err := med.categoryService.GetByIDs(categoryIDs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	categoryMap := make(map[string]domain.CategoryResponse)
// 	for _, ca := range categories {
// 		categoryMap[ca.ID] = *ca
// 	}

// 	var responseVDO []*domain.VideoResponse

// 	for _, result := range video {

// 		var categoryResponse domain.CategoryResponse
// 		if result.CategoryID != "" {
// 			if cat, ok := categoryMap[result.CategoryID]; ok {
// 				categoryResponse = cat
// 			}
// 		}

// 		resp := &domain.VideoResponse{
// 			Title:    result.Title,
// 			Content:  result.Content,
// 			URL:      result.VideoURL,
// 			Category: categoryResponse,
// 		}

// 		responseVDO = append(responseVDO, resp)

// 	}

// 	return responseVDO, nil
// }

// func (med *MediaService) GetShortVideoHome() ([]*domain.ShortVideo, error) {

// 	categories, err := med.categoryRepo.GetByName("Short Video")
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching category: %w", err)
// 	}

// 	short, err := med.MediaRepo.GetShortVideoHome(categories.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responseShort []*domain.ShortVideo

// 	for _, result := range short {

// 		resp := &domain.ShortVideo{
// 			Title: result.Title,
// 			URL:   result.VideoURL,
// 		}

// 		responseShort = append(responseShort, resp)

// 	}

// 	return responseShort, nil
// }
