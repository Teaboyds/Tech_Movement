package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"strings"
)

type BannerService struct {
	bannerRepo      port.BannerRepository
	cateRepo        port.CategoryRepository
	uploadSer       port.UploadService
	uploadRepo      port.UploadRepository
	CategoryService port.CategoryService
}

func NewBannerService(bannerRepo port.BannerRepository, cateRepo port.CategoryRepository, uploadSer port.UploadService, uploadRepo port.UploadRepository, CategoryService port.CategoryService) port.BannerService {
	return &BannerService{bannerRepo: bannerRepo, uploadSer: uploadSer, uploadRepo: uploadRepo, CategoryService: CategoryService, cateRepo: cateRepo}
}

func (ban *BannerService) CreateBanner( /*parameter*/ banner *domain.Banner) error {

	fmt.Printf("banner: %v\n", banner)

	if strings.TrimSpace(banner.Title) == "" || strings.TrimSpace(banner.ContentType) == "" {
		return fmt.Errorf("please input error")
	}

	_, err := ban.cateRepo.GetByID(banner.Category)
	if err != nil {
		return err
	}

	_, err = ban.uploadSer.ValidateImageIDs(banner.Img)
	if err != nil {
		return err
	}

	err = ban.bannerRepo.SaveBanner(banner)
	if err != nil {
		return err
	}

	fmt.Printf("banner.Status: %T\n banner.ContentTyep: %v\n", banner.Status, banner.ContentType)

	return nil
}

func (ban *BannerService) GetBanner(id string) (*domain.BannerClient, error) {

	banner, err := ban.bannerRepo.Retrive(id)
	if err != nil {
		return nil, err
	}

	cate, err := ban.cateRepo.GetByID(banner.Category)
	if err != nil {
		return nil, err
	}

	blueJeans := &domain.CategoryResponse{
		ID:   cate.ID,
		Name: cate.Name,
	}

	resp := &domain.BannerClient{
		Title:       banner.Title,
		ContentType: banner.ContentType,
		Status:      banner.Status,
		Category:    *blueJeans,
	}

	return resp, nil
}

func (ban *BannerService) GetBanners() ([]*domain.BannerClient, error) {
	banners, err := ban.bannerRepo.Retrives()
	if err != nil {
		return nil, err
	}
	imageIDBan := make(map[string]struct{})
	categoryIDBan := make(map[string]struct{})

	for _, banner := range banners {
		for _, imgID := range banner.Img {
			imageIDBan[imgID] = struct{}{}
		}
		if banner.Category != "" {
			categoryIDBan[banner.Category] = struct{}{}
		}
	}
	imageIDs := keysFromMap(imageIDBan)
	categories, err := ban.CategoryService.GetByIDs(keysFromMap(categoryIDBan))
	if err != nil {
		return nil, err
	}

	uploadFildes, err := ban.uploadRepo.GetFilesByIDs(imageIDs)
	if err != nil {
		return nil, err
	}

	uploadFileMap := make(map[string]domain.UploadFileResponse)
	for _, f := range uploadFildes {
		uploadFileMap[f.ID] = domain.UploadFileResponse{
			ID:       f.ID,
			Path:     f.Path,
			FileType: f.FileType,
		}
	}

	categoryMap := make(map[string]domain.CategoryResponse)
	for _, ca := range categories {
		categoryMap[ca.ID] = *ca
	}
	var responseBanners []*domain.BannerClient
	for _, result := range banners {
		var categoryResponse domain.CategoryResponse
		if result.Category != "" {
			if cat, ok := categoryMap[result.Category]; ok {
				categoryResponse = cat
			}
		}

		var images []domain.UploadFileResponse
		for _, imgID := range result.Img {
			if img, ok := uploadFileMap[imgID]; ok {
				images = append(images, img)
			}
		}

		resp := &domain.BannerClient{
			Title:       result.Title,
			ContentType: result.ContentType,
			Status:      result.Status,
			Category:    categoryResponse,
			Images:      images,
		}

		responseBanners = append(responseBanners, resp)
	}
	for _, item := range responseBanners {
		for j, img := range item.Images {
			item.Images[j].Path = utils.AttachBaseURLToImage(img.FileType, img.Path)
		}
	}

	return responseBanners, nil
}

func (ban *BannerService) Updated(id string, banner *domain.Banner) error {

	existingBanner, err := ban.bannerRepo.Retrive(id)
	if err != nil {
		return err
	}
	if banner.Category != "" {
		_, err := ban.cateRepo.GetByID(banner.Category)
		if err != nil {
			return err
		}
		existingBanner.Category = banner.Category
	} else {
		banner.Category = existingBanner.Category
	}
	if len(banner.Img) > 0 {
		_, err := ban.uploadSer.ValidateImageIDs(banner.Img)
		if err != nil {
			return err
		}
		existingBanner.Img = banner.Img
	} else {
		banner.Img = existingBanner.Img
	}

	fmt.Println("existingBanner.Img: ", existingBanner.Img)

	if banner.Title == "" {
		banner.Title = existingBanner.Title
	} else {
		existingBanner.Title = banner.Title
	}

	if banner.ContentType == "" {
		banner.ContentType = existingBanner.ContentType
	} else {
		existingBanner.ContentType = banner.ContentType
	}

	if !banner.Status {
		banner.Status = existingBanner.Status
	} else {
		existingBanner.Status = banner.Status
	}
	err = ban.bannerRepo.Updated(id, banner)
	if err != nil {
		return err
	}
	return err
}

func (ban *BannerService) Delete(id string) error {

	if err := ban.bannerRepo.Delete(id); err != nil {
		return err
	}
	return nil
}
