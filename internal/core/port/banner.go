package port

import "backend_tech_movement_hex/internal/core/domain"

type BannerRepository interface {
	SaveBanner( /*parameter*/ banner *domain.Banner) error
	Retrive(id string) (*domain.Banner, error)
	Retrives() ([]*domain.Banner, error)
	Updated(id string, banner *domain.Banner) error
	Delete(id string) error
}

type BannerService interface {
	CreateBanner( /*parameter*/ banner *domain.Banner) error
	GetBanner(id string) (*domain.BannerClient, error)
	GetBanners() ([]*domain.BannerClient, error)
	Updated(id string, banner *domain.Banner) error
	Delete(id string) error
}
