package port

import "backend_tech_movement_hex/internal/core/domain"

type BannerRepository interface {
	SaveBanner(banner *domain.Banner) error
	Retrive(id string) (*domain.Banner, error)
	Retrives(page_type string) ([]*domain.Banner, error)
	// Delete *ให้ลบภาพออกจาก localstorage ด้วย
	// Update *ทำแบบ ถ้าเปลี่ยนแค่บางฟิลด์ ฟิลด์อื่นไม่ต้องเปลี่ยน
	// Action *เปิดปิด action
	SaveBannerV2(banner *domain.BannerV2) error
	RetriveV2(id string) (*domain.BannerV2, error)
}

type BannerService interface {
	CreateBanner(banner *domain.Banner) error
	GetBanner(id string) (*domain.Banner, error)
	GetBanners(page_type string) ([]*domain.Banner, error)
	// Delete *ให้ลบภาพออกจาก localstorage ด้วย
	// Update *ทำแบบ ถ้าเปลี่ยนแค่บางฟิลด์ ฟิลด์อื่นไม่ต้องเปลี่ยน
	// Action *เปิดปิด action
	CreateBannerV2(banner *domain.BannerV2) error
	GetBannerV2(id string) (*domain.BannerV2, error)
}
