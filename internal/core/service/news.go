package service

import (
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
)

type NewsServiceImpl struct {
	newRepo      port.NewsRepository
	categoryRepo port.CategoryRepository
	fileRepo     port.UploadRepository
	cache        port.CacheRepository
}

func NewsService(newRepo port.NewsRepository, categoryRepo port.CategoryRepository, cache port.CacheRepository, fileRepo port.UploadRepository) port.NewsService {
	return &NewsServiceImpl{
		newRepo:      newRepo,
		categoryRepo: categoryRepo,
		cache:        cache,
		fileRepo:     fileRepo,
	}
}

func (n *NewsServiceImpl) CreateNews(news *d.NewsRequest) error {

	files, err := n.fileRepo.ValidateImageIDs(news.Image)
	if err != nil {
		return fmt.Errorf("failed to get image files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no image files found")
	}

	if len(files) > 3 {
		return fmt.Errorf("cannot attach more than 3 images")
	}

	err = n.newRepo.SaveNews(news)
	if err != nil {
		return err
	}

	return nil
}

// // / Get area ///

func (n *NewsServiceImpl) GetNewsByID(id string) (*d.NewsResponse, error) {

	news, err := n.newRepo.GetNewsByID(id)
	if err != nil {
		return nil, err
	}

	for i, image := range news.Image {
		news.Image[i].Path = utils.AttachBaseURLToImage(image.FileType, image.Path)
	}

	return news, nil
}

func (n *NewsServiceImpl) GetLastNews() ([]*d.HomePageLastedNewResponse, error) {

	news, err := n.newRepo.GetLastNews()
	if err != nil {
		return nil, err
	}

	for _, item := range news {
		for j, img := range item.Image {
			item.Image[j].Path = utils.AttachBaseURLToImage(img.Filetype, img.Path)
		}
	}

	return news, nil
}

func (n *NewsServiceImpl) GetTechnologyNews() ([]*d.HomePageLastedNewResponse, error) {
	news, err := n.newRepo.GetTechnologyNews()
	if err != nil {
		return nil, err
	}

	for _, item := range news {
		for j, img := range item.Image {
			item.Image[j].Path = utils.AttachBaseURLToImage(img.Filetype, img.Path)
		}
	}

	return news, nil
}

// func (n *NewsServiceImpl) GetNewsByCategoryHomePage(categoryID string) ([]d.News, error) {

// 	ctx, cancel := utils.NewTimeoutContext()
// 	defer cancel()

// 	ObjID, err := primitive.ObjectIDFromHex(categoryID)
// 	if err != nil {
// 		log.Println("Invalid category ID:", err)
// 		return nil, fmt.Errorf("invalid category ID format")
// 	}

// 	cacheKey := "news:category:" + ObjID.String()
// 	fmt.Printf("cacheKey: %v\n", cacheKey)

// 	var cacheNewsCategory []d.News
// 	err = n.cache.Get(ctx, cacheKey, &cacheNewsCategory)
// 	if err == nil && len(cacheNewsCategory) > 0 {
// 		log.Println("Cache Hit:", cacheKey)
// 		return cacheNewsCategory, nil
// 	}

// 	log.Println("Cache Miss:", cacheKey)

// 	news, err := n.repo.GetNewsByCategoryHomePage(categoryID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for i := range news {
// 		utils.AttachBaseURLToImage(&news[i])
// 		news[i].CreatedAtText = utils.ConvertTimeResponse(news[i].CreatedAt)
// 	}

// 	err = n.cache.Set(ctx, cacheKey, news, 5*time.Minute)
// 	if err != nil {
// 		log.Printf("Error setting cache for category %s: %v", cacheKey, err)
// 		return nil, err
// 	}

// 	return news, nil
// }

// func (n *NewsServiceImpl) GetNewsByWeek() ([]d.News, error) {

// 	news, err := n.repo.GetNewsByWeek()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for i := range news {
// 		utils.AttachBaseURLToImage(&news[i])
// 		news[i].CreatedAtText = utils.ConvertTimeResponse(news[i].CreatedAt)
// 	}

// 	return news, nil
// }

// /// Get area ///_

// func (n *NewsServiceImpl) UpdateNews(id string, req *d.UpdateNewsRequestResponse, filename string) error {

// 	existingNews, err := n.repo.GetNewsByID(id)
// 	if err != nil {
// 		return err
// 	}

// 	// debug ต้องนานดันไป get id cache มาหมดเวลาไป 5 ชม. บ่ได้หยัง //
// 	fmt.Println("req", filename)
// 	fmt.Println("existing", existingNews.Image)

// 	if req.Title != "" {
// 		existingNews.Title = req.Title
// 	}

// 	if req.Abstract != "" {
// 		existingNews.Abstract = req.Abstract
// 	}

// 	if req.Detail != "" {
// 		existingNews.Detail = req.Detail
// 	}

// 	oldImg := existingNews.Image
// 	if filename != "" {
// 		existingNews.Image = filename
// 	}

// 	if !utils.IsValidContentStatus(req.ContentStatus) {
// 		return fmt.Errorf("invalid content status")
// 	}

// 	if !utils.IsContentType(req.ContentType) {
// 		return fmt.Errorf("invalid content status")
// 	}

// 	if req.Category != "" {
// 		cat, err := n.categoryRepo.GetByID(req.Category)
// 		if err != nil {
// 			return err
// 		}
// 		existingNews.CategoryID = cat
// 	}

// 	existingNews.ContentStatus = req.ContentStatus
// 	existingNews.ContentType = req.ContentType

// 	if filename != "" && oldImg != "" && filename != oldImg {
// 		go func(oldImg string) {
// 			if err := n.repo.DeleteImg(oldImg); err != nil {
// 				log.Println("Failed to delete old image during update:", err)
// 			} else {
// 				log.Println("Old image deleted:", oldImg)
// 			}
// 		}(oldImg)
// 	}

// 	return n.repo.UpdateNews(id, existingNews)
// }

// func (n *NewsServiceImpl) Delete(id string) error {

// 	news, err := n.repo.GetNewsByID(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := n.repo.Delete(id); err != nil {
// 		return err
// 	}

// 	if news.Image != "" {
// 		err := n.repo.DeleteImg(news.Image)
// 		if err != nil {
// 			log.Println("Failed to delete image:", err)
// 		}
// 	}

// 	return nil
// }
