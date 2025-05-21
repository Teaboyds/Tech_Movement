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

	if !banner.Status.Home && !banner.Status.Media && !banner.Status.News && !banner.Status.Infographic {
		banner.Status.Home = true
		banner.Status.Media = true
		banner.Status.News = true
		banner.Status.Infographic = true
	}

	normalizeImage(&banner.DesktopImage)
	normalizeImage(&banner.MobileImage)

	if banner.Action == "" {
		banner.Action = "On"
	} else {
		actionCheck := utils.IsValidAction(banner.Action)
		if !actionCheck {
			return fmt.Errorf("action must be 'on' or 'off'")
		}
	}

	err := ban.bannerRepo.SaveBanner(banner)
	if err != nil {
		return err
	}

	return nil
}

func (ban *BannerService) GetBanner(id string) (*domain.Banner, error) {

	bannerResp, err := ban.bannerRepo.Retrive(id)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive banner in service : %v", err)
	}

	bannerResp.DesktopImage.Path = utils.AttachBaseURLToImageFolder("banner", bannerResp.DesktopImage.FileType, bannerResp.DesktopImage.Path)
	bannerResp.MobileImage.Path = utils.AttachBaseURLToImageFolder("banner", bannerResp.MobileImage.FileType, bannerResp.MobileImage.Path)

	return bannerResp, err
}

func (ban *BannerService) GetBanners(page_type string) ([]*domain.Banner, error) {

	bannerResp, err := ban.bannerRepo.Retrives(page_type)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch banners in service layer : %v", err)
	}

	for _, banners := range bannerResp {
		banners.DesktopImage.Path = utils.AttachBaseURLToImageFolder("banner", banners.DesktopImage.FileType, banners.DesktopImage.Path)
		banners.MobileImage.Path = utils.AttachBaseURLToImageFolder("banner", banners.MobileImage.FileType, banners.MobileImage.Path)
	}

	return bannerResp, err
}

func (ban *BannerService) CreateBannerV2(banner *domain.BannerV2) error {

	fmt.Printf("banner: %v\n", banner)

	if !banner.Status.Home && !banner.Status.Media && !banner.Status.News && !banner.Status.Infographic {
		banner.Status.Home = true
		banner.Status.Media = true
		banner.Status.News = true
		banner.Status.Infographic = true
	}

	if banner.Action == "" {
		banner.Action = "On"
	} else {
		actionCheck := utils.IsValidAction(banner.Action)
		if !actionCheck {
			return fmt.Errorf("action must be 'on' or 'off'")
		}
	}

	err := ban.bannerRepo.SaveBannerV2(banner)
	if err != nil {
		return err
	}

	return nil
}

func (ban *BannerService) GetBannerV2(id string) (*domain.BannerV2, error) {
	bannerResp, err := ban.bannerRepo.RetriveV2(id)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive banner in service : %v", err)
	}

	return bannerResp, err
}

// helper //
func normalizeImage(img *domain.ImageInfo) {
	img.Type = strings.TrimPrefix(img.Type, ".")
	img.Name = strings.ReplaceAll(img.Name, " ", "_")
	img.FileType = strings.TrimPrefix(img.FileType, "banner/")
}
