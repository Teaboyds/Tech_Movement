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

func (cats *CategoryService) Create(category *domain.CategoryRequest) error {

	existingName, err := cats.CatRepo.ExistsByName(category.Name)
	if err != nil {
		return err
	}

	if existingName {
		return fmt.Errorf("category name already exists")
	}

	err = cats.CatRepo.Create(category)
	if err != nil {
		return err
	}

	return nil
}

func (cats *CategoryService) GetByID(id string) (*domain.CategoryResponse, error) {

	category, err := cats.CatRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (cats *CategoryService) GetAll() ([]domain.CategoryResponse, error) {

	categories, err := cats.CatRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (cats *CategoryService) GetByName(name string) (*domain.Category, error) {
	return cats.CatRepo.GetByName(name)
}

func (cats *CategoryService) UpdateCategory(id string, category *domain.Category) error {

	if category.Name == "" {
		return fmt.Errorf("please input your new name")
	}

	err := cats.CatRepo.UpdateCategory(id, category)
	if err != nil {
		return err
	}

	return nil
}

func (cats *CategoryService) DeleteCategory(id string) error {
	return cats.CatRepo.DeleteCategory(id)
}
