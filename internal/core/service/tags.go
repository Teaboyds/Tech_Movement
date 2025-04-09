package service

import (
	r "backend_tech_movement_hex/internal/adapter/storage/mongodb/repository"
	dt "backend_tech_movement_hex/internal/core/domain"
	"errors"
	"fmt"
	"strings"
)

type TagsService struct {
	Tagrepo r.MongoTagsRepository
}

func TagsServiceImpl(Tagrepo r.MongoTagsRepository) *TagsService {
	return &TagsService{Tagrepo: Tagrepo}
}

func (tr *TagsService) CreateTags(tags dt.Tags) error {

	existingName, err := tr.Tagrepo.ExistsByName(tags.Name)
	if err != nil {
		return err
	}

	if existingName {
		return fmt.Errorf("category name already exists")
	}

	return tr.Tagrepo.SavaTags(tags)
}

var ErrTagIdNotProvided = errors.New("tag id is required")

func (tr *TagsService) GetTagsById(id string) (*dt.Tags, error) {

	if strings.TrimSpace(id) == "" {
		return nil, ErrTagIdNotProvided
	}

	return tr.Tagrepo.GetTagsById(id)
}
