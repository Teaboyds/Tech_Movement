package port

import (
	dif "backend_tech_movement_hex/internal/core/domain"
)

type InfographicRepository interface {
	CreateInfo(info *dif.Infographic) error
	GetInfoHome() ([]*dif.Infographic, error)
}

type InfographicService interface {
	CreateInfo(info *dif.InfographicRequest) error
	GetInfoHome() ([]*dif.InfographicRespose, error)
}
