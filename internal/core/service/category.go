package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
)

type CategoryService struct {
	CatRepo port.CategoryRepository
}

func NewCategoryService(CatRepo port.CategoryRepository) *CategoryService {
	return &CategoryService{CatRepo: CatRepo}
}

func (cats *CategoryService) Create(category *domain.Category) error {
	return cats.CatRepo.Create(category)
}

func (cats *CategoryService) GetByID(id string) (*domain.Category, error) {
	return cats.CatRepo.GetByID(id)
}

func (cats *CategoryService) GetAll() ([]domain.Category, error) {
	return cats.CatRepo.GetAll()
}

func (cats *CategoryService) GetByName(name string) (*domain.Category, error) {
	return cats.CatRepo.GetByName(name)
}
