package port

import (
	d "backend_tech_movement_hex/internal/core/domain"
)

type NewsRepository interface {
	SaveNews(news *d.News) error
	GetNewsByID(id string) (*d.News, error)
	GetLastNews() ([]*d.News, error)
	GetTechnologyNews() ([]*d.News, error)
	// GetNewsByTags(name string) ([]d.News, error)
	UpdateNews(id string, news *d.News) error
	Delete(id string) error
	EnsureNewsIndexs() error
	Find(catID, ConType, Sort string, limit, page int64) ([]*d.News, error)
}

type NewsService interface {
	CreateNews(news *d.News) error
	GetNewsByID(id string) (*d.NewsResponse, error)
	GetLastNews() ([]*d.NewsResponse, error)
	GetTechnologyNews() ([]*d.NewsResponse, error)
	UpdateNews(id string, req *d.News) error
	Delete(id string) error
	Find(catID, ConType, Sort, limit, page string) ([]*d.NewsResponse, error)
}
