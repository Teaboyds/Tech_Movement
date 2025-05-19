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
	DeleteMany(id []string) error
	EnsureNewsIndexs() error
	Find(catID, ConType, Sort, status, view string, limit, page int64) ([]*d.News, error)
	Count(catID, ConType, Status string) (int64, error)
}

type NewsService interface {
	CreateNews(news *d.News) error
	GetNewsByID(id string) (*d.NewsResponse, error)
	GetLastNews() ([]*d.NewsResponseV2, error)
	GetTechnologyNews() ([]*d.NewsResponseV2, error)
	UpdateNews(id string, req *d.News) error
	Delete(id string) error
	DeleteMany(id []string) error
	Find(catID, ConType, Sort, limit, page, status, view string) ([]*d.NewsResponseV2, error)
	Count(catID, ConType, Status, limit, page string) (*d.PaginationResp, error)
}
