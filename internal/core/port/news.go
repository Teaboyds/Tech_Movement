package port

import (
	d "backend_tech_movement_hex/internal/core/domain"
)

type NewsRepository interface {
	Create(news *d.News) error
	GetNewsPagination(lastID string, limit int) ([]d.News, error)
	GetNewsByID(id string) (*d.News, error)
	GetNewsByCategory(CategoryId string) ([]d.News, error)
	GetNewsByTags(name string) ([]d.News, error)
	UpdateNews(id string, news *d.News) error
	Delete(id string) error
	DeleteImg(path string) error
}

type NewsService interface {
	Create(news *d.News) error
	GetNewsPagination(lastID string, limit int) ([]d.News, error)
	GetNewsByID(id string) (*d.News, error)
	GetNewsByCategory(CategoryId string) ([]d.News, error)
	GetNewsByTags(name string) ([]d.News, error)
	UpdateNews(id string, news *d.News) error
	Delete(id string) error
}
