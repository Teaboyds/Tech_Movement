package port

import "backend_tech_movement_hex/internal/core/domain"

type BannerRepository interface {
	SaveBanner( /*parameter*/ banner *domain.Banner) error
	Retrive(id string) (*domain.Banner, error)
}

type BannerService interface {
	CreateBanner( /*parameter*/ banner *domain.Banner) error
	GetBanner(id string) (*domain.BannerClient, error)
}
