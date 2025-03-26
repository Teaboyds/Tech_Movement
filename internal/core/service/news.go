package service

import (
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
)

type NewsServiceImpl struct {
	repo         port.NewsRepository
	categoryRepo port.CategoryRepository
}

func NewsService(repo port.NewsRepository, categoryRepo port.CategoryRepository) port.NewsService {
	return &NewsServiceImpl{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (n *NewsServiceImpl) Create(news *d.News) error {

	// ถ้าไม่มี input เข้ามาเป็น cat id  ให้ default เป็น Uncategorized || *fallback value* //
	if news.CategoryID == nil {
		category, err := n.categoryRepo.GetByName("Uncategorized")
		if err != nil {
			defaultCategory := d.Category{
				Name: "Uncategorized",
			}
			err := n.categoryRepo.Create(&defaultCategory)
			if err != nil {
				return err
			}
			news.CategoryID = &defaultCategory
		} else {
			news.CategoryID = category
		}
	}

	// ถ้าไม่มี input เข้ามาเ  ให้ default เป็น Untagged //
	if len(news.Tag) == 0 {
		news.Tag = append(news.Tag, "Untagged")
	}

	news.CreatedAt, news.UpdatedAt = utils.SetTimestamps()

	return n.repo.Create(news)
}

func (n *NewsServiceImpl) GetNewsPagination(page int, limit int) ([]d.News, int, error) {
	return n.repo.GetNewsPagination(page, limit)
}

func (n *NewsServiceImpl) GetNewsByID(id string) (*d.News, error) {
	return n.repo.GetNewsByID(id)
}

func (n *NewsServiceImpl) UpdateNews(id string, news *d.News) error {
	return n.repo.UpdateNews(id, news)
}

func (n *NewsServiceImpl) Delete(id string) error {
	return n.repo.Delete(id)
}
