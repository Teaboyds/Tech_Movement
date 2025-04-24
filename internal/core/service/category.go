package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"fmt"
)

type CategoryService struct {
	CatRepo port.CategoryRepository
}

func NewCategoryService(CatRepo port.CategoryRepository) port.CategoryService {
	return &CategoryService{CatRepo: CatRepo}
}

func (cats *CategoryService) Create(category *domain.Category) error {

	existingName, err := cats.CatRepo.ExistsByName(category.Name)
	if err != nil {
		return err
	}

	if existingName {
		return fmt.Errorf("category name already exists")
	}

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

func (cats *CategoryService) UpdateCategory(id string, category *domain.Category) error {
	return cats.CatRepo.UpdateCategory(id, category)
}

func (cats *CategoryService) DeleteCategory(id string) error {
	return nil
}
