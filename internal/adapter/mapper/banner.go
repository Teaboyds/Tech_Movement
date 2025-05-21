package mapper

import "backend_tech_movement_hex/internal/core/domain"

func BannerRequestToDomain(req domain.BannerRequestV2, status domain.StatusType) *domain.BannerV2 {
	return &domain.BannerV2{
		DesktopImage: req.DesktopImage,
		MobileImage:  req.MobileImage,
		Status: domain.StatusType{
			Home:        status.Home,
			Media:       status.Media,
			News:        status.News,
			Infographic: status.Infographic,
		},
		LinkUrl:   req.LinkUrl,
		Action:    req.Action,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
}

func BannerDomainToResponse(req domain.BannerV2, img domain.MetaData, imgMo domain.MetaData, status domain.StatusType) *domain.BannerResponseV2 {
	return &domain.BannerResponseV2{
		ID: req.ID,
		DesktopImageUrl: domain.MetaData{
			Alt:      img.Alt,
			Url:      img.Url,
			Size:     img.Size,
			MimeType: img.MimeType,
			Type:     img.Type,
		},
		MobileImageUrl: domain.MetaData{
			Alt:      imgMo.Alt,
			Url:      imgMo.Url,
			Size:     imgMo.Size,
			MimeType: imgMo.MimeType,
			Type:     imgMo.Type,
		},
		Status: domain.StatusType{
			Home:        status.Home,
			Media:       status.Media,
			News:        status.News,
			Infographic: status.Infographic,
		},
		Action:    req.Action,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
}
