package port

import dt "backend_tech_movement_hex/internal/core/domain"

type TagsRepository interface {
	SavaTags(tags dt.Tags) error
	GetTagsById(id string) (*dt.Tags, error)
	GetTagsAll() ([]dt.Tags, error)
	EditTags(id string, tags dt.Tags) error
	DeleteTags(id string) error
	ExistsByName(name string) (bool, error) // interface repo หาเพื่อเอามาเช็คใน service ว่าสร้างไปหรือยังออะ //
}

type TagsService interface {
	CreateTags(tags dt.Tags) error
	GetTagsById(id string) (*dt.Tags, error)
	GetTagsAll() ([]dt.Tags, error)
	UpdateTags(id string, tags dt.Tags) error
	DeleteTags(id string) error
}
