package port

import (
	"backend_tech_movement_hex/internal/core/domain"
)

type MediaRepository interface {
	CreateMedia(media *domain.Media) error
	GetVideoHome(cateId string) ([]*domain.Media, error)
	GetShortVideoHome(cateId string) ([]*domain.Media, error)
}

type MediaService interface {
	CreateMedia(media *domain.Media) error
	GetVideoHome() ([]*domain.VideoResponse, error)
	GetShortVideoHome() ([]*domain.ShortVideo, error)
}
