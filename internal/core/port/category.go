package port

import (
	dc "backend_tech_movement_hex/internal/core/domain"
)

type CategoryRepository interface {
	Create(category *dc.Category) error
	GetByID(id string) (*dc.Category, error)
	GetByName(name string) (*dc.Category, error)
	GetAll() ([]dc.Category, error)
}

type CategoryService interface {
	Create(category *dc.Category) error
	GetByID(id string) (*dc.Category, error)
	GetByName(name string) (*dc.Category, error)
	GetAll() ([]dc.Category, error)
}
