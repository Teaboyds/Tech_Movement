package port

import (
	d "backend_tech_movement_hex/internal/core/domain"
)

type NewsRepository interface {
	Create(news *d.News) error
	GetNewsPagination(page int, limit int) ([]d.News, int, error)
	GetNewsByID(id string) (*d.News, error)
	UpdateNews(id string, news *d.News) error
	Delete(id string) error
}

type NewsService interface {
	Create(news *d.News) error
	GetNewsPagination(page int, limit int) ([]d.News, int, error)
	GetNewsByID(id string) (*d.News, error)
	UpdateNews(id string, news *d.News) error
	Delete(id string) error
}
