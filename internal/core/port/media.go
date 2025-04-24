package port

import "backend_tech_movement_hex/internal/core/domain"

type MediaRepository interface {
	SaveMedia(media *domain.Media) error
	EnsureMediaIndexs() error
	GetLastMedia() ([]domain.Media, error)
}

type MediaService interface {
	CreateMedia(media *domain.Media) error
	GetLastMedia() ([]domain.Media, error)
}
