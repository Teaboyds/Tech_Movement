// primary adapters //
package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"context"
	"fmt"
	"log"
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

type ImageDataResponse struct {
	ImagePath string `json:"ImagePath"`
	ImageName string `json:"imageName"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
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

	title := c.FormValue("title")
	detail := c.FormValue("detail")
	statusStr := c.FormValue("status")
	contentStatus := c.FormValue("content_status")
	categoryID := c.FormValue("category")
	tag := c.FormValue("tag")

	if !utils.IsValidContentStatus(contentStatus) {
		log.Printf("Invalid content status received: %s", contentStatus)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Content Status"})
	}

	status := statusStr == "true"

	tagList := strings.Split(tag, ",")

	fileName, err := utils.UploadFile(c, "image", 5*1024*1024, "./upload/image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// ดึงไอดีของ cat มาเพื่อจะได้นำมาใส่ตอน create //
	category, _ := h.CategoryService.GetByID(categoryID)

	newNews := domain.News{
		Title:         title,
		Detail:        detail,
		Image:         fileName,
		CategoryID:    category,
		Tag:           tagList,
		Status:        status,
		ContentStatus: contentStatus,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	if err := h.service.Create(&newNews); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.cacheService.DeletePattern(context.Background(), "News_Keys_*"); err != nil {
		log.Printf("Error deleting cache key: %v", err)
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

	CategoryID := c.Params("id")

	fmt.Println(CategoryID)

	newsList, err := h.service.GetNewsByCategory(CategoryID)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	if len(newsList) == 0 {
		fmt.Println("No news found for CategoryID:", CategoryID)
	}

	// instance ดึงตัว name category เพื่อนำไปแสดงใน response //
	categoryName := newsList[0].CategoryID.Name

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"category": categoryName,
		"news":     newsList,
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

	// ดึงข่าวมาเปรียบเทียบ //
	existingNews, err := h.service.GetNewsByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve the news. Please try again later.",
		})
	}
	oldImg := existingNews.Image

	title := c.FormValue("title")
	detail := c.FormValue("detail")
	statusStr := c.FormValue("status")
	contentStatus := c.FormValue("content_status")
	categoryID := c.FormValue("category")
	tag := c.FormValue("tag")

	if !utils.IsValidContentStatus(contentStatus) {
		log.Printf("Invalid content status: %s", contentStatus)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Content Status"})
	}

	status := statusStr == "true"
	tagList := strings.Split(tag, ",")

	newfile := oldImg
	if file, err := c.FormFile("image"); err == nil && file != nil {
		newfile, err = utils.UploadFile(c, "image", 5*1024*1024, "./upload/image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if categoryID != "" {
		category, err := h.CategoryService.GetByID(categoryID)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrResponse{
				Error: "Failed to retrieve category",
			})
		}
		existingNews.CategoryID = category
	}

	if title != "" {
		existingNews.Title = title
	}
	if detail != "" {
		existingNews.Detail = detail
	}
	if tag != "" {
		existingNews.Tag = tagList
	}

	existingNews.Status = status
	existingNews.ContentStatus = contentStatus
	existingNews.Image = newfile
	existingNews.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := h.service.UpdateNews(id, existingNews); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot Update News cause Internal Server Error ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "News Updated Successfully",
		"data":    existingNews,
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
