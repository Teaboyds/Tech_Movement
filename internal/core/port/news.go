package port

import d "backend_tech_movement_hex/internal/core/domain"

type NewsRepository interface {
	SaveNews(news *d.NewsRequest) error
	GetNewsByID(id string) (*d.NewsResponse, error)
	GetLastNews() ([]*d.HomePageLastedNewResponse, error)
	GetTechnologyNews() ([]*d.HomePageLastedNewResponse, error)
	// GetNewsByCategoryHomePage(categoryID string) ([]d.News, error)
	// GetNewsByTags(name string) ([]d.News, error)
	// UpdateNews(id string, news *d.News) error
	// Delete(id string) error
	// DeleteImg(path string) error
	EnsureNewsIndexs() error
}

type NewsService interface {
	CreateNews(news *d.NewsRequest) error
	GetNewsByID(id string) (*d.NewsResponse, error)
	GetLastNews() ([]*d.HomePageLastedNewResponse, error)
	GetTechnologyNews() ([]*d.HomePageLastedNewResponse, error)
	// GetNewsByCategoryHomePage(categoryID string) ([]d.News, error)
	// UpdateNews(id string, req *d.UpdateNewsRequestResponse, filename string) error
	// Delete(id string) error
}
