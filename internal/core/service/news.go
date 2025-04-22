package service

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/repository"
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"context"
	"fmt"
	"log"
	"time"
)

type NewsServiceImpl struct {
	repo         *repository.MongoNewsRepository
	categoryRepo port.CategoryRepository
	cache        port.CacheRepository
}

func NewsService(repo *repository.MongoNewsRepository, categoryRepo port.CategoryRepository, cache port.CacheRepository) port.NewsService {
	return &NewsServiceImpl{
		repo:         repo,
		categoryRepo: categoryRepo,
		cache:        cache,
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

	return n.repo.Create(news)
}

// / Get area ///
func (n *NewsServiceImpl) GetNewsPagination(lastID string, limit int) ([]d.News, error) {
	return n.repo.GetNewsPagination(lastID, limit)
}

func (n *NewsServiceImpl) GetNewsByID(id string) (*d.News, error) {

	news, err := n.repo.GetNewsByID(id)
	if err != nil {
		return nil, err
	}

	utils.AttachBaseURLToImage(news)

	return news, nil
}

func (n *NewsServiceImpl) GetNewsByCategory(categoryID string, lastID string) ([]d.News, string, error) {

	news, nextCursor, err := n.repo.GetNewsByCategory(categoryID, lastID)
	if err != nil {
		return nil, "", err
	}

	for i := range news {
		utils.AttachBaseURLToImage(&news[i])
	}

	return news, nextCursor, nil
}

func (n *NewsServiceImpl) GetLastNews() ([]d.News, error) {

	ctx := context.Background()

	versionKey := "news:latest:version"
	cachePrefix := "news:latest:data:"

	var version string
	err := n.cache.Get(ctx, versionKey, &version)
	if err != nil || version == "" {
		version = "v1"
	}

	cacheKey := cachePrefix + version

	var cacheNews []d.News
	err = n.cache.Get(ctx, cacheKey, &cacheNews)
	if err == nil && len(cacheNews) > 0 {
		log.Println("Cache Hit:", cacheKey)
		return cacheNews, nil
	}

	log.Println("Cache Miss:", cacheKey)

	news, err := n.repo.GetLastNews()
	if err != nil {
		return nil, err
	}

	for i := range news {
		utils.AttachBaseURLToImage(&news[i])
	}

	_ = n.cache.Set(ctx, cacheKey, news, 3*time.Minute)

	return news, nil
}

/// Get area ///

func (n *NewsServiceImpl) UpdateNews(id string, req *d.UpdateNewsRequestResponse, filename string) error {

	existingNews, err := n.repo.GetNewsByID(id)
	if err != nil {
		return err
	}

	// debug ต้องนานดันไป get id cache มาหมดเวลาไป 5 ชม. บ่ได้หยัง //
	fmt.Println("req", filename)
	fmt.Println("existing", existingNews.Image)

	if req.Title != "" {
		existingNews.Title = req.Title
	}

	if req.Detail != "" {
		existingNews.Detail = req.Detail
	}

	oldImg := existingNews.Image
	if filename != "" {
		existingNews.Image = filename
	}

	if !utils.IsValidContentStatus(req.ContentStatus) {
		return fmt.Errorf("invalid content status")
	}

	if !utils.IsContentType(req.ContentType) {
		return fmt.Errorf("invalid content status")
	}

	if req.Category != "" {
		cat, err := n.categoryRepo.GetByID(req.Category)
		if err != nil {
			return err
		}
		existingNews.CategoryID = cat
	}

	existingNews.ContentStatus = req.ContentStatus
	existingNews.ContentType = req.ContentType

	if filename != "" && oldImg != "" && filename != oldImg {
		go func(oldImg string) {
			if err := n.repo.DeleteImg(oldImg); err != nil {
				log.Println("Failed to delete old image during update:", err)
			} else {
				log.Println("Old image deleted:", oldImg)
			}
		}(oldImg)
	}

	return n.repo.UpdateNews(id, existingNews)
}

func (n *NewsServiceImpl) Delete(id string) error {

	news, err := n.repo.GetNewsByID(id)
	if err != nil {
		return err
	}

	if err := n.repo.Delete(id); err != nil {
		return err
	}

	if news.Image != "" {
		err := n.repo.DeleteImg(news.Image)
		if err != nil {
			log.Println("Failed to delete image:", err)
		}
	}

	return nil
}
