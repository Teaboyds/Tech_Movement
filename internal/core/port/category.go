package port

import (
	dc "backend_tech_movement_hex/internal/core/domain"
)

type CategoryRepository interface {
	Create(category *dc.CategoryRequest) error
	GetByID(id string) (*dc.CategoryResponse, error)
	GetByName(name string) (*dc.Category, error)
	GetAll() ([]dc.CategoryResponse, error)
	GetByIDs(ids []string) ([]*dc.CategoryResponse, error)
	UpdateCategory(id string, category *dc.Category) error
	DeleteCategory(id string) error
	ExistsByName(name string) (bool, error)
}

type CategoryService interface {
	Create(category *dc.CategoryRequest) error
	GetByID(id string) (*dc.CategoryResponse, error)
	GetByName(name string) (*dc.Category, error)
	GetAll() ([]dc.CategoryResponse, error)
	UpdateCategory(id string, category *dc.Category) error
	DeleteCategory(id string) error
}
