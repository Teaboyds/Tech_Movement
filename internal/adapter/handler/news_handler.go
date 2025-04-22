// primary adapters //
package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsHandler struct {
	service         port.NewsService
	CategoryService port.CategoryRepository
	cacheService    port.CacheRepository
}

func NewNewsHandler(
	service port.NewsService,
	CategoryService port.CategoryRepository,
	cacheService port.CacheRepository,
) *NewsHandler {
	return &NewsHandler{
		service:         service,
		CategoryService: CategoryService,
		cacheService:    cacheService,
	}
}

// CreateNews godoc
// @Summary Create a New News
// @Description api สร้างหน้าข่าวใหม่
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param news body domain.News true "News Data"
// @Success 201 {object} domain.News
// @Failure 404 {object} domain.ErrResponse
// @Failure 500 {object} domain.ErrResponse
// @Router /news [post]
func (h *NewsHandler) CreateNews(c *fiber.Ctx) error {

	var input domain.NewsRequest

	if err := c.BodyParser(&input); err != nil {
		log.Printf("news bad request %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid News Reuquest",
		})
	}

	if err := utils.ValidateNewsInput(&input); err != nil {
		log.Printf("news bad validator request %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Validator News Reuquest",
		})
	}

	status := input.Status == "true"
	tagList := strings.Split(input.Tag, ",")

	fileName, err := utils.UploadFile(c, "image", 5*1024*1024, "./upload/image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// ดึงไอดีของ cat มาเพื่อจะได้นำมาใส่ตอน create //
	category, _ := h.CategoryService.GetByID(input.Category)

	newNews := domain.News{
		Title:         input.Title,
		Detail:        input.Detail,
		Image:         fileName,
		CategoryID:    category,
		Tag:           tagList,
		Status:        status,
		ContentStatus: input.ContentStatus,
		ContentType:   input.ContentType,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	if err := h.service.Create(&newNews); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(newNews)
}

// GetNewsByPage godoc
// @Summary Get News Pagination
// @Description ดึงข้อมูลข่าวแบบ pagination โดยดึงจาก lastid กล่าวคือหากใส่ id ของข่าวแล้วหลังจาก id ลงไปนั้นตาม limit
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param lastID query string false "ID of the last news item from the previous page"
// @Param limit query int false "Number of news per page (default: 10)"
// @Success 200 {object} domain.News
// @Failure 400 {object} domain.ErrResponse "Invalid request parameters"
// @Failure 500 {object} domain.ErrResponse "Cannot Fetch Data."
// @Router /news [get]
func (h *NewsHandler) GetNewsByPage(c *fiber.Ctx) error {

	lastID := c.Query("lastID")
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	news, err := h.service.GetNewsPagination(lastID, limit)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: "Cannot Fetch Data.",
		})
	}

	return c.JSON(fiber.Map{
		"data":   news,
		"limit":  limit,
		"lastID": news[len(news)-1].ID,
	})
}

func (h *NewsHandler) GetNewsByID(c *fiber.Ctx) error {
	id := c.Params("id")

	news, err := h.service.GetNewsByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(domain.ErrResponse{
				Error: "News Not Found",
			})
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: "Cannot Fetch News Data.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": news,
	})
}

func (h *NewsHandler) GetNewsByCategory(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	lastID := c.Query("lastID") // ปรับจาก Params → Query parameter

	newsResult, nextCursor, err := h.service.GetNewsByCategory(categoryID, lastID)
	if err != nil {
		fmt.Println("GetNewsByCategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot Fetch News By Category",
		})
	}

	if len(newsResult) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"category": "",
			"news":     []domain.NewsHomePageResponse{},
			"next":     "",
		})
	}

	var lastedNews []domain.NewsHomePageResponse
	for _, news := range newsResult {
		lastedNews = append(lastedNews, domain.NewsHomePageResponse{
			Title:     news.Title,
			Detail:    news.Detail,
			Image:     news.Image,
			Category:  news.CategoryID.Name,
			CreatedAt: news.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"category": newsResult[0].CategoryID.Name,
		"news":     lastedNews,
		"next":     nextCursor,
	})
}

func (h *NewsHandler) GetLastNews(c *fiber.Ctx) error {

	lastNews, err := h.service.GetLastNews()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth category data",
		})
	}

	var lastedNews []domain.HomePageLastedNewResponse
	for _, news := range lastNews {
		lastedNews = append(lastedNews, domain.HomePageLastedNewResponse{
			Title:     news.Title,
			Detail:    news.Detail,
			Image:     news.Image,
			Category:  news.CategoryID.Name,
			CreatedAt: news.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ข่าวล่าสุด",
		"data":    lastedNews,
	})
}

func (h *NewsHandler) UpdateNews(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		log.Println("Invalid ObjectID format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ObjectID",
		})
	}

	var news domain.UpdateNewsRequestResponse
	if err := c.BodyParser(&news); err != nil {
		log.Printf("news bad request %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid News Reuquest",
		})
	}

	fileName, err := utils.UploadFile(c, "image", 5*1024*1024, "./upload/image")
	if err != nil && err != http.ErrMissingFile {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.UpdateNews(id, &news, fileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot Update News cause Internal Server Error ",
		})
	}

	newNews := domain.UpdateNewsRequestResponse{
		Title:         news.Title,
		Detail:        news.Detail,
		Image:         fileName,
		Category:      news.Category,
		Tag:           news.Tag,
		Status:        news.Status,
		ContentStatus: news.ContentStatus,
		ContentType:   news.ContentType,
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "News Updated Successfully",
		"data":    newNews,
	})
}

func (h *NewsHandler) DeleteNews(c *fiber.Ctx) error {

	id := c.Params("id")
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ObjectID",
		})
	}

	err = h.service.Delete(id)
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
