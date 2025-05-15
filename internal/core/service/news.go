package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type NewsServiceImpl struct {
	newRepo         port.NewsRepository
	categoryRepo    port.CategoryRepository
	categoryService port.CategoryService
	fileRepo        port.UploadRepository
	cache           port.CacheRepository
}

func NewsService(newRepo port.NewsRepository, categoryRepo port.CategoryRepository, cache port.CacheRepository, fileRepo port.UploadRepository, categoryService port.CategoryService) port.NewsService {
	return &NewsServiceImpl{
		newRepo:         newRepo,
		categoryRepo:    categoryRepo,
		cache:           cache,
		fileRepo:        fileRepo,
		categoryService: categoryService,
	}
}

func (n *NewsServiceImpl) CreateNews(news *d.News) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	_, err := n.fileRepo.ValidateImageIDs(news.Image)
	if err != nil {
		return fmt.Errorf("failed to get image files: %s", err)
	}

	_, err = n.categoryRepo.GetByID(news.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to get category: %s ", err)
	}

	if len(news.Image) == 0 {
		return fmt.Errorf("no image files found")
	}

	if len(news.Image) > 3 {
		return fmt.Errorf("cannot attach more than 3 images")
	}

	err = n.newRepo.SaveNews(news)
	if err != nil {
		return err
	}

	findCache := "news_cache:" + news.CategoryID + "*"
	err = n.cache.DeletePattern(ctx, findCache)
	if err != nil {
		return err
	} else {
		log.Printf("Deleted cache with pattern: %s", findCache)
	}

	return nil
}

// // / Get area ///

func (n *NewsServiceImpl) GetNewsByID(id string) (*d.NewsResponse, error) {

	news, err := n.newRepo.GetNewsByID(id)
	if err != nil {
		return nil, err
	}

	category, err := n.categoryService.GetByID(news.CategoryID)
	if err != nil {
		return nil, err
	}

	// ดึงข้อมูลไฟล์จาก FileRepo (ดึงข้อมูลไฟล์ที่มี ID)
	uploadFiles, err := n.fileRepo.GetFilesByIDs(news.Image)
	if err != nil {
		return nil, err
	}

	var images []domain.UploadFileResponseHomePage
	for _, file := range uploadFiles {
		images = append(images, domain.UploadFileResponseHomePage{
			ID:       file.ID,
			Path:     file.Path,
			Filetype: file.FileType,
		})
	}

	response := &domain.NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		Image:       images,    // เติมข้อมูลภาพที่แปลงจาก ID
		CategoryID:  *category, // เติมข้อมูล Category ที่แปลงจาก ID
		Tag:         news.Tag,
		Status:      news.Status,
		ContentType: news.ContentType,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	}

	for i, img := range response.Image {
		response.Image[i].Path = utils.AttachBaseURLToImage(img.Filetype, img.Path)
	}

	return response, nil
}

func (n *NewsServiceImpl) GetLastNews() ([]*d.NewsResponse, error) {

	lastedNews, err := n.newRepo.GetLastNews()
	if err != nil {
		return nil, err
	}

	imageIDMap := make(map[string]struct{})
	categoryIDMap := make(map[string]struct{})

	for _, news := range lastedNews {
		for _, imgID := range news.Image {
			imageIDMap[imgID] = struct{}{}
		}
		if news.CategoryID != "" {
			categoryIDMap[news.CategoryID] = struct{}{}
		}
	}

	imageIDs := keysFromMap(imageIDMap)
	categoryIDs := keysFromMap(categoryIDMap)

	uploadFiles, err := n.fileRepo.GetFilesByIDs(imageIDs)
	if err != nil {
		return nil, err
	}

	categories, err := n.categoryService.GetByIDs(categoryIDs)
	if err != nil {
		return nil, err
	}

	uploadFileMap := make(map[string]domain.UploadFileResponseHomePage)
	for _, f := range uploadFiles {
		uploadFileMap[f.ID] = domain.UploadFileResponseHomePage{
			ID:       f.ID,
			Path:     f.Path,
			Filetype: f.FileType,
		}
	}

	categoryMap := make(map[string]domain.CategoryResponse)
	for _, ca := range categories {
		categoryMap[ca.ID] = *ca
	}

	var responseNews []*domain.NewsResponse

	for _, news := range lastedNews {
		var images []domain.UploadFileResponseHomePage
		for _, imgID := range news.Image {
			if img, ok := uploadFileMap[imgID]; ok {
				images = append(images, img)
			}
		}

		var categoryResponse domain.CategoryResponse
		if news.CategoryID != "" {
			if cat, ok := categoryMap[news.CategoryID]; ok {
				categoryResponse = cat
			}
		}

		resp := &domain.NewsResponse{
			ID:          news.ID,
			Title:       news.Title,
			Description: news.Description,
			Content:     news.Content,
			Image:       images,
			CategoryID:  categoryResponse,
			Tag:         news.Tag,
			Status:      news.Status,
			ContentType: news.ContentType,
			CreatedAt:   news.CreatedAt,
			UpdatedAt:   news.UpdatedAt,
		}
		responseNews = append(responseNews, resp)
	}

	for _, item := range responseNews {
		for j, img := range item.Image {
			item.Image[j].Path = utils.AttachBaseURLToImage(img.Filetype, img.Path)
		}
	}

	return responseNews, nil
}

func (n *NewsServiceImpl) GetTechnologyNews() ([]*d.NewsResponse, error) {
	news, err := n.newRepo.GetTechnologyNews()
	if err != nil {
		return nil, err
	}

	imageIDMap := make(map[string]struct{})
	categoryIDMap := make(map[string]struct{})

	for _, news := range news {
		for _, imgID := range news.Image {
			imageIDMap[imgID] = struct{}{}
		}
		if news.CategoryID != "" {
			categoryIDMap[news.CategoryID] = struct{}{}
		}
	}

	imageIDs := keysFromMap(imageIDMap)
	categoryIDs := keysFromMap(categoryIDMap)

	uploadFiles, err := n.fileRepo.GetFilesByIDs(imageIDs)
	if err != nil {
		return nil, err
	}

	categories, err := n.categoryService.GetByIDs(categoryIDs)
	if err != nil {
		return nil, err
	}

	uploadFileMap := make(map[string]domain.UploadFileResponseHomePage)
	for _, f := range uploadFiles {
		uploadFileMap[f.ID] = domain.UploadFileResponseHomePage{
			ID:       f.ID,
			Path:     f.Path,
			Filetype: f.FileType,
		}
	}

	categoryMap := make(map[string]domain.CategoryResponse)
	for _, ca := range categories {
		categoryMap[ca.ID] = *ca
	}

	var responseNews []*domain.NewsResponse

	for _, news := range news {
		var images []domain.UploadFileResponseHomePage
		for _, imgID := range news.Image {
			if img, ok := uploadFileMap[imgID]; ok {
				images = append(images, img)
			}
		}

		var categoryResponse domain.CategoryResponse
		if news.CategoryID != "" {
			if cat, ok := categoryMap[news.CategoryID]; ok {
				categoryResponse = cat
			}
		}

		fmt.Printf("categoryResponse: %v\n", categoryResponse.CategoryType)

		resp := &domain.NewsResponse{
			ID:          news.ID,
			Title:       news.Title,
			Description: news.Description,
			Content:     news.Content,
			Image:       images,
			CategoryID:  categoryResponse,
			Tag:         news.Tag,
			Status:      news.Status,
			ContentType: news.ContentType,
			CreatedAt:   news.CreatedAt,
			UpdatedAt:   news.UpdatedAt,
		}
		responseNews = append(responseNews, resp)
	}

	for _, item := range responseNews {
		for j, img := range item.Image {
			item.Image[j].Path = utils.AttachBaseURLToImage(img.Filetype, img.Path)
		}
	}

	return responseNews, nil
}

func (n *NewsServiceImpl) Find(catID, ConType, Sort, limit, page string) ([]*d.NewsResponse, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	Sort = strings.ToLower(Sort)
	if Sort != "asc" && Sort != "desc" {
		return nil, fmt.Errorf("sort must be 'asc' or 'desc'")
	}

	fmt.Printf("page: %v\n", page)

	parseLimit, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, err
	}

	parsePage, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return nil, err
	}

	CacheKey := "news_cache:" + catID + "_" + ConType + "_" + page + "_" + limit
	if ConType == "" {
		CacheKey = "news_cache:" + catID + "_" + "All" + "_" + page + "_" + limit
	}
	// cache area //
	if page == "1" {
		var cacheNewsCategory []*d.NewsResponse
		err = n.cache.Get(ctx, CacheKey, &cacheNewsCategory)
		if err == nil && len(cacheNewsCategory) > 0 {
			fmt.Printf("cacheKey: %v\n", CacheKey)
			log.Println("Cache Hit:", CacheKey)
			return cacheNewsCategory, nil
		}
		log.Println("Cache Miss:", CacheKey)
	}

	findingNemo, err := n.newRepo.Find(catID, ConType, Sort, parseLimit, parsePage)
	if err != nil {
		return nil, err
	}

	imageIDMap := make(map[string]struct{})
	categoryIDMap := make(map[string]struct{})

	for _, news := range findingNemo {
		for _, imgID := range news.Image {
			imageIDMap[imgID] = struct{}{}
		}
		if news.CategoryID != "" {
			categoryIDMap[news.CategoryID] = struct{}{}
		}
	}

	imageIDs := keysFromMap(imageIDMap)
	categoryIDs := keysFromMap(categoryIDMap)

	uploadFiles, err := n.fileRepo.GetFilesByIDs(imageIDs)
	if err != nil {
		return nil, err
	}

	categories, err := n.categoryService.GetByIDs(categoryIDs)
	if err != nil {
		return nil, err
	}

	uploadFileMap := make(map[string]domain.UploadFileResponseHomePage)
	for _, f := range uploadFiles {
		uploadFileMap[f.ID] = domain.UploadFileResponseHomePage{
			ID:       f.ID,
			Path:     f.Path,
			Filetype: f.FileType,
		}
	}

	categoryMap := make(map[string]domain.CategoryResponse)
	for _, ca := range categories {
		categoryMap[ca.ID] = *ca
	}

	var responseNews []*domain.NewsResponse

	for _, news := range findingNemo {
		var images []domain.UploadFileResponseHomePage
		for _, imgID := range news.Image {
			if img, ok := uploadFileMap[imgID]; ok {
				images = append(images, img)
			}
		}

		var categoryResponse domain.CategoryResponse
		if news.CategoryID != "" {
			if cat, ok := categoryMap[news.CategoryID]; ok {
				categoryResponse = cat
			}
		}

		resp := &domain.NewsResponse{
			ID:          news.ID,
			Title:       news.Title,
			Description: news.Description,
			Content:     news.Content,
			Image:       images,
			CategoryID:  categoryResponse,
			Tag:         news.Tag,
			Status:      news.Status,
			ContentType: news.ContentType,
			CreatedAt:   news.CreatedAt,
			UpdatedAt:   news.UpdatedAt,
		}
		responseNews = append(responseNews, resp)
	}

	for _, item := range responseNews {
		for j, img := range item.Image {
			item.Image[j].Path = utils.AttachBaseURLToImage(img.Filetype, img.Path)
		}
	}

	if page == "1" {
		err = n.cache.Set(ctx, CacheKey, responseNews, 5*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("error setting cache for category %s: %v", CacheKey, err)
		}
	}

	return responseNews, nil
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

// var cacheNewsCategory []d.News
// err = n.cache.Get(ctx, cacheKey, &cacheNewsCategory)
// if err == nil && len(cacheNewsCategory) > 0 {
// 	log.Println("Cache Hit:", cacheKey)
// 	return cacheNewsCategory, nil
// }

// 	log.Println("Cache Miss:", cacheKey)

// 	news, err := n.repo.GetNewsByCategoryHomePage(categoryID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for i := range news {
// 		utils.AttachBaseURLToImage(&news[i])
// 		news[i].CreatedAtText = utils.ConvertTimeResponse(news[i].CreatedAt)
// 	}

// err = n.cache.Set(ctx, cacheKey, news, 5*time.Minute)
// if err != nil {
// 	log.Printf("Error setting cache for category %s: %v", cacheKey, err)
// 	return nil, err
// }

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

func (n *NewsServiceImpl) UpdateNews(id string, req *d.News) error {

	existingNews, err := n.newRepo.GetNewsByID(id)
	if err != nil {
		return err
	}

	if req.CategoryID != "" {
		_, err = n.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			return fmt.Errorf("failed to get category : %w", err)
		}
		existingNews.CategoryID = req.CategoryID
	} else {
		req.CategoryID = existingNews.CategoryID
	}

	if len(req.Image) > 0 {
		img, err := n.fileRepo.ValidateImageIDs(req.Image)
		if err != nil {
			return fmt.Errorf("failed to get image files: %w", err)
		}
		existingNews.Image = img
	} else {
		req.Image = existingNews.Image
	}

	// debug ต้องนานดันไป get id cache มาหมดเวลาไป 5 ชม. บ่ได้หยัง //
	fmt.Println("existing", existingNews.Image)

	if req.Title == "" {
		req.Title = existingNews.Title
	} else {
		existingNews.Title = req.Title
	}

	if req.Description == "" {
		req.Description = existingNews.Description
	} else {
		existingNews.Description = req.Description
	}

	if req.Content == "" {
		req.Content = existingNews.Content
	} else {
		existingNews.Content = req.Content
	}

	if len(req.Tag) == 0 {
		req.Tag = existingNews.Tag
	} else {
		existingNews.Tag = req.Tag
	}

	if !req.Status {
		req.Status = existingNews.Status
	} else {
		existingNews.Status = req.Status
	}

	if req.ContentType == "" {
		req.ContentType = existingNews.ContentType
	} else {
		existingNews.ContentType = req.ContentType
	}

	fmt.Printf("req: %v\n", req)

	err = n.newRepo.UpdateNews(id, req)
	if err != nil {
		return err
	}

	return err
}

func (n *NewsServiceImpl) Delete(id string) error {

	if err := n.newRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

// helper func //
func keysFromMap(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
