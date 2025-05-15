// primary adapters //
package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type NewsHandler struct {
	mediaService    port.MediaService
	service         port.NewsService
	categoryService port.CategoryService
	cacheService    port.CacheRepository
	infoGraphic     port.InfographicService
}

func NewNewsHandler(
	mediaService port.MediaService,
	service port.NewsService,
	CategoryService port.CategoryService,
	cacheService port.CacheRepository,
	infoGraphic port.InfographicService,
) *NewsHandler {
	return &NewsHandler{
		mediaService:    mediaService,
		service:         service,
		categoryService: CategoryService,
		cacheService:    cacheService,
		infoGraphic:     infoGraphic,
	}
}

func (h *NewsHandler) CreateNews(c *fiber.Ctx) error {
	var input domain.NewsRequest

	if err := c.BodyParser(&input); err != nil {
		log.Printf("news bad request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid News Request",
		})
	}

	if err := utils.FixMultipartArray(c, &input); err != nil {
		log.Printf("fix multipart array error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Form Array Field",
		})
	}

	fmt.Printf("input after fix: %+v\n", input)

	if err := utils.ValidateNewsInput(&input); err != nil {
		log.Printf("news bad validator request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Validator News Request",
		})
	}

	req := &domain.News{
		Title:       input.Title,
		Description: input.Description,
		Content:     input.Content,
		Image:       input.Image,
		CategoryID:  input.Category,
		Tag:         input.Tag,
		Status:      input.Status,
		ContentType: input.ContentType,
	}

	err := h.service.CreateNews(req)
	if err != nil {
		log.Printf("create news error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(input)
}

func (h *NewsHandler) GetNewsByID(c *fiber.Ctx) error {
	id := c.Params("id")

	news, err := h.service.GetNewsByID(id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrResponse{
			Error: "News Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": news,
	})
}

func (h *NewsHandler) Find(c *fiber.Ctx) error {

	find := c.Queries()
	catId := find["categoryID"]
	conType := find["contentType"]
	sort := find["sort"]
	limit := find["limit"]
	page := find["page"]

	findingNemo, err := h.service.Find(catId, conType, sort, limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": findingNemo,
	})
}

// get news by api //
func (h *NewsHandler) GetLastNews(c *fiber.Ctx) error {

	lastNews, err := h.service.GetLastNews()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth lastNews data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ข่าวล่าสุด",
		"data":    lastNews,
	})
}

func (h *NewsHandler) GetTechNews(c *fiber.Ctx) error {

	TechNews, err := h.service.GetTechnologyNews()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth lastNews data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ข่าว Technology",
		"data":    TechNews,
	})
}

func (h *NewsHandler) GetVideoHome(c *fiber.Ctx) error {

	VDO, err := h.mediaService.GetVideoHome()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth Video data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ข่าว Technology",
		"data":    VDO,
	})
}

func (h *NewsHandler) GetShortVideoHome(c *fiber.Ctx) error {

	ShortVDO, err := h.mediaService.GetShortVideoHome()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth ShortVideo data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ข่าว Technology",
		"data":    ShortVDO,
	})
}

func (h *NewsHandler) GetInfoHome(c *fiber.Ctx) error {

	Info, err := h.infoGraphic.GetInfoHome()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth Info data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Infographic",
		"data":    Info,
	})
}

func (h *NewsHandler) GetHomePage(c *fiber.Ctx) error {

	Video, err := h.mediaService.GetVideoHome()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth Video data",
		})
	}

	TechNews, err := h.service.GetTechnologyNews()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth TechNews data",
		})
	}

	lastNews, err := h.service.GetLastNews()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth lastNews data",
		})
	}

	shortVDO, err := h.mediaService.GetShortVideoHome()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth shortVdo data",
		})
	}

	Infographic, err := h.infoGraphic.GetInfoHome()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth infographic data",
		})
	}

	resp := domain.Home{
		Message:        "Home Landing Page",
		Video:          Video,
		TechnologyNews: TechNews,
		LastedNews:     lastNews,
		Short:          shortVDO,
		Infographic:    Infographic,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// func (h *NewsHandler) GetNewsByCategory(c *fiber.Ctx) error {

// 	id := c.Query("id", "")

// 	newsCat, err := h.service.GetNewsByCategoryHomePage(id)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "cannot fecth category data",
// 		})
// 	}

// 	var lastedNews []domain.NewsHomeCategoryPageResponse
// 	for _, news := range newsCat {
// 		lastedNews = append(lastedNews, domain.NewsHomeCategoryPageResponse{
// 			Title:     news.Title,
// 			Abstract:  news.Abstract,
// 			Detail:    news.Detail,
// 			Image:     news.Image,
// 			Category:  news.CategoryID.Name,
// 			CreatedAt: news.CreatedAtText,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "ข่าวล่าสุด",
// 		"data":    lastedNews,
// 	})
// }

func (h *NewsHandler) UpdateNews(c *fiber.Ctx) error {
	id := c.Params("id")

	var news domain.NewsRequest
	if err := c.BodyParser(&news); err != nil {
		log.Printf("news bad request %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid News Reuquest",
		})
	}

	newNews := domain.News{
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		Image:       news.Image,
		CategoryID:  news.Category,
		Tag:         news.Tag,
		Status:      news.Status,
		ContentType: news.ContentType,
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := h.service.UpdateNews(id, &newNews); err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot Update News cause Internal Server Error ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "News Updated Successfully",
		"data":    newNews,
	})
}

func (h *NewsHandler) DeleteNews(c *fiber.Ctx) error {

	id := c.Params("id")

	err := h.service.Delete(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot Delete News cause Internal Server Error ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "News Delete Successfully",
	})

}

func (h *NewsHandler) DeleteManyNews(c *fiber.Ctx) error {

	var news domain.DeleteManyID

	if err := c.BodyParser(&news); err != nil {
		log.Printf("news ID bad request %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid News ID Reuquest",
		})
	}

	fmt.Printf("news: %v\n", news)

	fieldMap := map[string]string{
		"ids": "IDs",
	}

	if err := utils.FixMultipartArrayV2(c, &news, fieldMap); err != nil {
		log.Printf("fix multipart array error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Form Array Field",
		})
	}

	if err := h.service.DeleteMany(news.IDs); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot Delete News cause Internal Server Error ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "News Delete Successfully",
		"news":    news.IDs,
	})
}
