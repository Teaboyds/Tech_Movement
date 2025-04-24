package service

import (
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsServiceImpl struct {
	repo         port.NewsRepository
	categoryRepo port.CategoryRepository
	cache        port.CacheRepository
}

func NewsService(repo port.NewsRepository, categoryRepo port.CategoryRepository, cache port.CacheRepository) port.NewsService {
	return &NewsServiceImpl{
		repo:         repo,
		categoryRepo: categoryRepo,
		cache:        cache,
	}
}

func (n *NewsServiceImpl) Create(news *d.News) error {

	ctx := context.Background()

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

	err := n.repo.Create(news)
	if err != nil {
		return err
	}

	cachePrefix := "news:latest:data:"
	cacheKey := cachePrefix + "latest"
	err = n.cache.Delete(ctx, cacheKey)
	if err != nil {
		return err
	}

	newsCatKeys := "news:category:" + news.CategoryID.ID.String()
	err = n.cache.Delete(ctx, newsCatKeys)
	if err != nil {
		return err
	}

	return nil
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

	cachePrefix := "news:latest:data:"
	cacheKey := cachePrefix + "latest"

	var cacheNews []d.News
	err := n.cache.Get(ctx, cacheKey, &cacheNews)
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

	ttl := 10 * time.Minute
	err = n.cache.Set(ctx, cacheKey, news, ttl)
	if err != nil {
		log.Printf("Error setting cache for %v: %v", cacheKey, err)
		return nil, err
	}

	return news, nil
}

func (n *NewsServiceImpl) GetNewsByCategoryHomePage(categoryID string) ([]d.News, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		log.Println("Invalid category ID:", err)
		return nil, fmt.Errorf("invalid category ID format")
	}

	cacheKey := "news:category:" + ObjID.String()
	fmt.Printf("cacheKey: %v\n", cacheKey)

	var cacheNewsCategory []d.News
	err = n.cache.Get(ctx, cacheKey, &cacheNewsCategory)
	if err == nil && len(cacheNewsCategory) > 0 {
		log.Println("Cache Hit:", cacheKey)
		return cacheNewsCategory, nil
	}

	log.Println("Cache Miss:", cacheKey)

	news, err := n.repo.GetNewsByCategoryHomePage(categoryID)
	if err != nil {
		return nil, err
	}

	for i := range news {
		utils.AttachBaseURLToImage(&news[i])
		news[i].CreatedAtText = utils.ConvertTimeResponse(news[i].CreatedAt)
	}

	err = n.cache.Set(ctx, cacheKey, news, 5*time.Minute)
	if err != nil {
		log.Printf("Error setting cache for category %s: %v", cacheKey, err)
		return nil, err
	}

	return news, nil
}

func (n *NewsServiceImpl) GetNewsByWeek() ([]d.News, error) {

	news, err := n.repo.GetNewsByWeek()
	if err != nil {
		return nil, err
	}

	for i := range news {
		utils.AttachBaseURLToImage(&news[i])
		news[i].CreatedAtText = utils.ConvertTimeResponse(news[i].CreatedAt)
	}

	return news, nil
}

/// Get area ///_

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

	if req.Abstract != "" {
		existingNews.Abstract = req.Abstract
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
