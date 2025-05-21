package mongoMapper

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	"backend_tech_movement_hex/internal/core/domain"
	"time"
)

func toDomainImage(img models.ImageInfo) domain.ImageInfo {
	return domain.ImageInfo{
		Path:     img.Path,
		FileType: img.FileType,
	}
}

func toDomainStatus(s models.StatusType) domain.StatusType {
	return domain.StatusType{
		Home:        s.Home,
		Media:       s.Media,
		News:        s.News,
		Infographic: s.Infographic,
	}
}

func BannerToDomain(banner models.MongoBanner) *domain.Banner {
	return &domain.Banner{
		ID:           banner.ID.Hex(),
		DesktopImage: toDomainImage(banner.DesktopImage),
		MobileImage:  toDomainImage(banner.MobileImage),
		Status:       toDomainStatus(banner.Status),
		LinkUrl:      banner.LinkUrl,
		Action:       banner.Action,
		CreatedAt:    banner.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    banner.UpdatedAt.Format(time.RFC3339),
	}
}

func MongoBannerToDomain(banner models.MongoBannerV2) *domain.BannerV2 {
	return &domain.BannerV2{
		ID:           banner.ID.Hex(),
		DesktopImage: banner.DesktopImage,
		MobileImage:  banner.MobileImage,
		Status:       toDomainStatus(banner.Status),
		LinkUrl:      banner.LinkUrl,
		Action:       banner.Action,
		CreatedAt:    banner.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    banner.UpdatedAt.Format(time.RFC3339),
	}
}
