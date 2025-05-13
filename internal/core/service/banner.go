package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"fmt"
	"strings"
)

type BannerService struct {
	bannerRepo port.BannerRepository
	cateRepo   port.CategoryRepository
	uploadSer  port.UploadService
}

func NewBannerService(bannerRepo port.BannerRepository, cateRepo port.CategoryRepository, uploadSer port.UploadService) port.BannerService {
	return &BannerService{bannerRepo: bannerRepo, cateRepo: cateRepo, uploadSer: uploadSer}
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
		Category:    blueJeans,
	}

	return resp, nil
}
