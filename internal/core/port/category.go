package port

import (
	dc "backend_tech_movement_hex/internal/core/domain"
)

type CategoryRepository interface {
	Create(category *dc.Category) error
	GetByID(id string) (*dc.Category, error)
	GetByName(name string) (*dc.Category, error)
	GetAll() ([]dc.Category, error)
	GetByIDs(ids []string) ([]*dc.Category, error)
	UpdateCategory(id string, category *dc.Category) error
	DeleteCategory(id string) error
	ExistsByName(name string) (bool, error)
}

type CategoryService interface {
	Create(category *dc.Category) error
	GetByID(id string) (*dc.Category, error)
	GetByIDs(ids []string) ([]*dc.CategoryResponse, error)
	GetByName(name string) (*dc.Category, error)
	GetAll() ([]dc.CategoryResponse, error)
	UpdateCategory(id string, category *dc.Category) error
	DeleteCategory(id string) error
}
