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

	err = cats.CatRepo.Create(category)
	if err != nil {
		return err
	}

	return nil
}

func (cats *CategoryService) GetByID(id string) (*domain.Category, error) {

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

	var cateDTO []domain.CategoryResponse
	for _, cate := range categories {
		cateDTO = append(cateDTO, domain.CategoryResponse{
			ID:           cate.ID,
			Name:         cate.Name,
			CategoryType: cate.CategoryType,
		})
	}

	return cateDTO, nil
}

func (cats *CategoryService) GetByIDs(ids []string) ([]*domain.Category, error) {

	Cate, err := cats.CatRepo.GetByIDs(ids)
	if err != nil {
		return nil, err
	}

	var responses []*domain.Category
	for _, category := range Cate {
		resp := &domain.Category{
			ID:           category.ID,
			Name:         category.Name,
			CategoryType: category.CategoryType,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
		}
		responses = append(responses, resp)
	}

	return responses, err
}

func (cats *CategoryService) GetByName(name string) (*domain.Category, error) {
	return cats.CatRepo.GetByName(name)
}

func (cats *CategoryService) UpdateCategory(id string, category *domain.Category) error {

	existingCate, err := cats.CatRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("category not found")
	}

	if category.Name != "" {
		existingCate.Name = category.Name
	}

	if category.CategoryType != "" {
		existingCate.CategoryType = category.CategoryType
	}

	err = cats.CatRepo.UpdateCategory(id, category)
	if err != nil {
		return err
	}

	return nil
}

func (cats *CategoryService) DeleteCategory(id string) error {

	err := cats.CatRepo.DeleteCategory(id)
	if err != nil {
		return err
	}

	return nil
}
