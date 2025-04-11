package service

import (
	dt "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"fmt"
)

type TagsServiceImpl struct {
	Tagrepo port.TagsRepository
}

func TagsService(Tagrepo port.TagsRepository) *TagsServiceImpl {
	return &TagsServiceImpl{Tagrepo: Tagrepo}
}

func (tr *TagsServiceImpl) CreateTags(tags *dt.Tags) error {

	if tags.Name == "" {
		return fmt.Errorf("tag name is required")
	}

	existingName, err := tr.Tagrepo.ExistsByName(tags.Name)
	if err != nil {
		return err
	}

	if existingName {
		return fmt.Errorf("category name already exists")
	}

	return tr.Tagrepo.SavaTags(tags)
}

func (tr *TagsServiceImpl) GetTagsByIdArray(id []string) ([]*dt.Tags, error) {
	return tr.Tagrepo.GetTagsByIdArray(id)
}

func (tr *TagsServiceImpl) GetTagsAll() ([]dt.Tags, error) {
	return tr.Tagrepo.GetAllTags()
}

func (tr *TagsServiceImpl) UpdateTags(id string, tags *dt.Tags) error {

	if tags.Name == "" {
		return fmt.Errorf("tag name is required")
	}

	return tr.Tagrepo.EditTags(id, tags)
}

func (tr *TagsServiceImpl) DeleteTags(id string) error {
	return tr.Tagrepo.DeleteTags(id)
}

func (tr *TagsServiceImpl) GetByID(id string) (*dt.Tags, error) {
	return tr.Tagrepo.GetByID(id)
}
