package service

import (
	"backend_tech_movement_hex/internal/core/domain"
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type NewsServiceImpl struct {
	newRepo         port.NewsRepository
	categoryRepo    port.CategoryRepository
	categoryService port.CategoryService
	fileRepo        port.UploadRepository
	fileService     port.UploadService
	cache           port.CacheRepository
}

func NewsService(newRepo port.NewsRepository, categoryRepo port.CategoryRepository, cache port.CacheRepository, fileRepo port.UploadRepository, categoryService port.CategoryService, fileService port.UploadService) port.NewsService {
	return &NewsServiceImpl{
		newRepo:         newRepo,
		categoryRepo:    categoryRepo,
		cache:           cache,
		fileRepo:        fileRepo,
		fileService:     fileService,
		categoryService: categoryService,
	}
}

func (n *NewsServiceImpl) CreateNews(news *d.News) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	_, err := n.fileRepo.GetFileByID(news.ThumnailID)
	if err != nil {
		return fmt.Errorf("failed to get thumnail files: %s", err)
	}

	_, err = n.fileRepo.ValidateImageIDs(news.ImageIDs)
	if err != nil {
		return fmt.Errorf("failed to get image files: %s", err)
	}

	_, err = n.categoryRepo.GetByID(news.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to get category: %s ", err)
	}

	if len(news.ImageIDs) == 0 {
		return fmt.Errorf("no image files found")
	}

	if len(news.ImageIDs) > 3 {
		return fmt.Errorf("cannot attach more than 3 images")
	}

	if news.View == "" {
		news.View = "0"
		_, err := strconv.Atoi(news.View)
		if err != nil {
			return fmt.Errorf("failed to parse View: %s", err)
		}
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

// Get news by id ไว้ใช้หน้าข่าว เดี่ยว ๆ //
func (n *NewsServiceImpl) GetNewsByID(id string) (*d.NewsResponse, error) {

	news, err := n.newRepo.GetNewsByID(id)
	if err != nil {
		return nil, err
	}

	cate, err := n.categoryService.GetByID(news.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("cannot fecth category data in Find func : %v", err)
	}

	categoryResponse := &domain.CategoryResponse{
		ID:           cate.ID,
		Name:         cate.Name,
		CategoryType: cate.CategoryType,
	}

	thumnail, err := n.fileService.GetFileByID(news.ThumnailID)
	if err != nil {
		return nil, fmt.Errorf("cannot fecth thumnail data in Find func : %v", err)
	}

	// ดึงข้อมูลไฟล์จาก FileRepo (ดึงข้อมูลไฟล์ที่มี ID)
	uploadFiles, err := n.fileService.GetFilesByIDsVTest(news.ImageIDs)
	if err != nil {
		return nil, fmt.Errorf("cannot fecth img data in Find func : %v", err)
	}

	// แปลง type
	var fileValues []domain.UploadFileResponse
	for _, f := range uploadFiles {
		if f != nil {
			fileValues = append(fileValues, domain.UploadFileResponse{
				ID:       f.ID,
				Name:     f.Name,
				Path:     f.Path,
				FileType: f.FileType,
				Type:     f.Type,
			})
		}
	}

	thumnailDto := &domain.UploadFileResponse{
		ID:       thumnail.ID,
		Name:     thumnail.Name,
		Path:     thumnail.Path,
		FileType: thumnail.FileType,
		Type:     thumnail.Type,
	}

	response := &domain.NewsResponse{
		ID:          news.ID,
		ThumnailID:  *thumnailDto,
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		ImageIDs:    fileValues,        // เติมข้อมูลภาพที่แปลงจาก ID
		CategoryID:  *categoryResponse, // เติมข้อมูล Category ที่แปลงจาก ID
		Tags:        news.Tags,
		Status:      news.Status,
		ContentType: news.ContentType,
		PageView:    news.View,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	}

	return response, nil
}

// last news หน้า home หรือ หน้า landing page news //
func (n *NewsServiceImpl) GetLastNews() ([]*d.NewsResponseV2, error) {
	lastedNews, err := n.newRepo.GetLastNews()
	if err != nil {
		return nil, err
	}

	categoryIDMap := make(map[string]struct{})
	thumnailIDMap := make(map[string]struct{})

	for _, news := range lastedNews {
		if news.CategoryID != "" {
			categoryIDMap[news.CategoryID] = struct{}{}
		}
		if news.CategoryID != "" {
			thumnailIDMap[news.ThumnailID] = struct{}{}
		}
	}

	categoryIDs := keysFromMap(categoryIDMap)
	thumnailIDs := keysFromMap(thumnailIDMap)

	thumnail, err := n.fileService.GetFilesByIDsVTest(thumnailIDs)
	if err != nil {
		return nil, err
	}

	categories, err := n.categoryService.GetByIDs(categoryIDs)
	if err != nil {
		return nil, err
	}

	thumnailMap := make(map[string]domain.UploadFileResponse)
	for _, t := range thumnail {
		thumnailMap[t.ID] = d.UploadFileResponse{
			ID:       t.ID,
			Name:     t.Name,
			Path:     t.Path,
			FileType: t.FileType,
			Type:     t.Type,
		}
	}

	categoryMap := make(map[string]domain.Category)
	for _, ca := range categories {
		categoryMap[ca.ID] = *ca
	}

	var responseNews []*domain.NewsResponseV2

	for _, news := range lastedNews {

		var categoryResponse domain.Category
		if news.CategoryID != "" {
			if cat, ok := categoryMap[news.CategoryID]; ok {
				categoryResponse = cat
			}
		}

		var thumnailResponse domain.UploadFileResponse
		if news.ThumnailID != "" {
			if thum, ok := thumnailMap[news.ThumnailID]; ok {
				thumnailResponse = thum
			}
		}

		resp := &domain.NewsResponseV2{
			ID:          news.ID,
			ThumnailID:  thumnailResponse,
			Title:       news.Title,
			Description: news.Description,
			Content:     news.Content,
			CategoryID:  categoryResponse,
			Tags:        news.Tags,
			Status:      news.Status,
			ContentType: news.ContentType,
			View:        news.View,
			CreatedAt:   news.CreatedAt,
		}
		responseNews = append(responseNews, resp)
	}

	return responseNews, nil
}

func (n *NewsServiceImpl) GetTechnologyNews() ([]*d.NewsResponseV2, error) {

	lastedNews, err := n.newRepo.GetTechnologyNews()
	if err != nil {
		return nil, err
	}

	categoryIDMap := make(map[string]struct{})
	thumnailIDMap := make(map[string]struct{})

	for _, news := range lastedNews {
		if news.CategoryID != "" {
			categoryIDMap[news.CategoryID] = struct{}{}
		}
		if news.CategoryID != "" {
			thumnailIDMap[news.ThumnailID] = struct{}{}
		}
	}

	categoryIDs := keysFromMap(categoryIDMap)
	thumnailIDs := keysFromMap(thumnailIDMap)

	thumnail, err := n.fileService.GetFilesByIDsVTest(thumnailIDs)
	if err != nil {
		return nil, err
	}

	categories, err := n.categoryService.GetByIDs(categoryIDs)
	if err != nil {
		return nil, err
	}

	thumnailMap := make(map[string]domain.UploadFileResponse)
	for _, t := range thumnail {
		thumnailMap[t.ID] = d.UploadFileResponse{
			ID:       t.ID,
			Name:     t.Name,
			Path:     t.Path,
			FileType: t.FileType,
			Type:     t.Type,
		}
	}

	categoryMap := make(map[string]domain.Category)
	for _, ca := range categories {
		categoryMap[ca.ID] = *ca
	}

	var responseNews []*domain.NewsResponseV2

	for _, news := range lastedNews {

		var categoryResponse domain.Category
		if news.CategoryID != "" {
			if cat, ok := categoryMap[news.CategoryID]; ok {
				categoryResponse = cat
			}
		}

		var thumnailResponse domain.UploadFileResponse
		if news.ThumnailID != "" {
			if thum, ok := thumnailMap[news.ThumnailID]; ok {
				thumnailResponse = thum
			}
		}

		resp := &domain.NewsResponseV2{
			ID:          news.ID,
			ThumnailID:  thumnailResponse,
			Title:       news.Title,
			Description: news.Description,
			Content:     news.Content,
			CategoryID:  categoryResponse,
			Tags:        news.Tags,
			Status:      news.Status,
			ContentType: news.ContentType,
			View:        news.View,
			CreatedAt:   news.CreatedAt,
		}
		responseNews = append(responseNews, resp)
	}

	return responseNews, nil
}

// เอาไว้ query หน้า
func (n *NewsServiceImpl) Find(catID, ConType, Sort, limit, page, status, view, search string) ([]*d.NewsResponseV2, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	fmt.Printf("limit: %v\n", limit)

	Sort = strings.ToLower(Sort)
	view = strings.ToLower(view)

	if Sort == "" {
	} else if Sort != "newest" && Sort != "oldest" {
		return nil, fmt.Errorf("sort must be 'newest' or 'oldest'")
	}

	if view == "" {
	} else if view != "asc" && view != "desc" {
		return nil, fmt.Errorf("sort must be 'asc' or 'desc'")
	}

	parseLimit, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse limit :%v", err)
	}

	parsePage, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse page :%v", err)
	}

	CacheKey := "news_cache:" + catID + "_" + ConType + "_" + page + "_" + limit + "_" + status + "_" + view
	if ConType == "" {
		CacheKey = "news_cache:" + catID + "_" + "All" + "_" + page + "_" + limit + "_" + status
	}

	// cache area //
	if page == "1" && search == "" {
		var cacheNewsCategory []*d.NewsResponseV2
		err = n.cache.Get(ctx, CacheKey, &cacheNewsCategory)
		if err == nil && len(cacheNewsCategory) > 0 {
			fmt.Printf("cacheKey: %v\n", CacheKey)
			log.Println("Cache Hit:", CacheKey)
			return cacheNewsCategory, nil
		}
		log.Println("Cache Miss:", CacheKey)
	}

	findingNemo, err := n.newRepo.Find(catID, ConType, Sort, status, view, search, parseLimit, parsePage)
	if err != nil {
		return nil, err
	}

	thumbIDs := make([]string, 0)
	catIDs := make([]string, 0)
	for _, n := range findingNemo {
		if n.ThumnailID != "" {
			thumbIDs = append(thumbIDs, n.ThumnailID)
		}
		if n.CategoryID != "" {
			catIDs = append(catIDs, n.CategoryID)
		}
	}

	thumbnails, _ := n.fileService.GetFilesByIDsVTest(thumbIDs)
	categories, _ := n.categoryService.GetByIDs(catIDs)

	thumbMap := make(map[string]domain.UploadFileResponse)
	for _, t := range thumbnails {
		thumbMap[t.ID] = domain.UploadFileResponse{
			ID:       t.ID,
			Name:     t.Name,
			Path:     t.Path,
			FileType: t.FileType,
			Type:     t.Type,
		}
	}

	catMap := make(map[string]domain.Category)
	for _, c := range categories {
		catMap[c.ID] = *c
	}

	var result []*domain.NewsResponseV2
	for _, nemo := range findingNemo {
		resp := &domain.NewsResponseV2{
			ID:          nemo.ID,
			Title:       nemo.Title,
			Description: nemo.Description,
			Content:     nemo.Content,
			Status:      nemo.Status,
			ContentType: nemo.ContentType,
			CreatedAt:   nemo.CreatedAt,
			Tags:        nemo.Tags,
			View:        nemo.View,
		}

		if thumb, ok := thumbMap[nemo.ThumnailID]; ok {
			resp.ThumnailID = thumb
		}
		if cat, ok := catMap[nemo.CategoryID]; ok {
			resp.CategoryID = cat
		}
		result = append(result, resp)
	}

	if page == "1" && search == "" {
		err = n.cache.Set(ctx, CacheKey, result, 5*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("error setting cache for category %s: %v", CacheKey, err)
		}
	}

	return result, nil
}

func (n *NewsServiceImpl) Count(catID, ConType, Status, limit, page string) (*d.PaginationResp, error) {

	total, err := n.newRepo.Count(catID, ConType, Status)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch count data:%s", err)
	}

	if limit == "" {
		limit = "10"
	}

	parseLimit, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse limit :%s", err)
	}

	totalPageSum := int64(math.Ceil(float64(total) / float64(parseLimit)))

	paginate := &domain.PaginationResp{
		TotalItems:  strconv.FormatInt(total, 10),
		TotalPages:  strconv.FormatInt(totalPageSum, 10),
		CurrentPage: page,
		PageSize:    limit,
	}

	return paginate, err
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

	if len(req.ImageIDs) > 0 {
		img, err := n.fileRepo.ValidateImageIDs(req.ImageIDs)
		if err != nil {
			return fmt.Errorf("failed to get image files: %w", err)
		}
		existingNews.ImageIDs = img
	} else {
		req.ImageIDs = existingNews.ImageIDs
	}

	// debug ต้องนานดันไป get id cache มาหมดเวลาไป 5 ชม. บ่ได้หยัง //
	fmt.Println("existing", existingNews.ImageIDs)

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

	if len(req.Tags) == 0 {
		req.Tags = existingNews.Tags
	} else {
		existingNews.Tags = req.Tags
	}

	if req.Status == "" {
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

func (n *NewsServiceImpl) DeleteMany(id []string) error {

	if err := n.newRepo.DeleteMany(id); err != nil {
		return fmt.Errorf("cannot delete News: %s", err)
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
