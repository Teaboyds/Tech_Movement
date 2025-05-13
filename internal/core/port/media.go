package port

import (
	"backend_tech_movement_hex/internal/core/domain"
)

type MediaRepository interface {
	CreateMedia(media *domain.MediaRequest) error
	GetVideoHome(cateId string) ([]*domain.VideoResponse, error)
	GetShortVideoHome() ([]*domain.ShortVideo, error)
}

type MediaService interface {
	CreateMedia(media *domain.MediaRequest) error
	GetVideoHome() ([]*domain.VideoResponse, error)
	GetShortVideoHome() ([]*domain.ShortVideo, error)
}
