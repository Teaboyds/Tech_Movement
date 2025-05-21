package port

import (
	"backend_tech_movement_hex/internal/core/domain"
	dif "backend_tech_movement_hex/internal/core/domain"
)

type InfographicRepository interface {
	CreateInfo(info *dif.Infographic) error
	// GetInfoHome() ([]*dif.Infographic, error)
	Retrive(id string) (*dif.Infographic, error)
	RetrivesInfographic(cateId, sort, view, limit, page string) ([]*domain.Infographic, error)
}

type InfographicService interface {
	CreateInfo(info *dif.Infographic) error
	// GetInfoHome() ([]*dif.InfographicRespose, error)
	GetInfographic(id string) (*dif.Infographic, error)
	GetInfographics(cateId, sort, view, limit, page string) ([]*domain.Infographic, error)
}
