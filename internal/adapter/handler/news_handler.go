// primary adapters //
package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsHandler struct {
	service         port.NewsService
	CategoryService port.CategoryRepository
}

func NewNewsHandler(service port.NewsService, CategoryService port.CategoryRepository) *NewsHandler {
	return &NewsHandler{service: service, CategoryService: CategoryService}
}

// CreateNews godoc
// @Summary Create a news article
// @Description Create a new news article
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param news body domain.News true "News Data"
// @Success 201 {object} domain.News
// @Failure 404 {object} domain.ErrResponse
// @Failure 500 {object} domain.ErrResponse
// @Router / [post]
func (h *NewsHandler) CreateNews(c *fiber.Ctx) error {
	var news UpdateNewsRequest
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid News input"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot open file",
		})
	}
	defer fileContent.Close()

	data, err := io.ReadAll(fileContent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot read file content",
		})
	}

	news.Image = base64.StdEncoding.EncodeToString(data)

	// ดึงไอดีของ category มาเพื่อจะได้นำมาใส่ตอน create //
	category, _ := h.CategoryService.GetByID(news.Category)

	newNews := domain.News{
		Title:      news.Title,
		Detail:     news.Detail,
		Image:      news.Image,
		CategoryID: category,
		Tag:        news.Tag,
	}

	if err := h.service.Create(&newNews); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(newNews)
}

// GetNewsByPage godoc
// @Summary Get News Pagination
// @Description Get paginated news by specifying page and limit
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of news per page (default: 10)"
// @Success 200 {object} domain.News
// @Failure 400 {object} domain.ErrResponse "Invalid request parameters"
// @Failure 500 {object} domain.ErrResponse "Cannot Fetch Data."
// @Router / [get]
func (h *NewsHandler) GetNewsByPage(c *fiber.Ctx) error {

	lastID := c.Query("lastID")
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	fmt.Println(lastID)

	news, err := h.service.GetNewsPagination(lastID, limit)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: "Cannot Fetch Data.",
		})
	}

	for i := range news {
		if err := utils.ProcessImageToURL(&news[i]); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrResponse{
				Error: "Failed to process image.",
			})
		}
	}

	return c.JSON(fiber.Map{
		"data":  news,
		"limit": limit,
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
			Error: "Cannot Fetch Dataa.",
		})
	}

	if err := utils.ProcessImageToURL(news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: "Failed to process image.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": news,
	})
}

type UpdateNewsRequest struct {
	Title     string   `json:"title"`
	Detail    string   `json:"detail"`
	Image     string   `json:"image"`
	Category  string   `json:"category"`
	Tag       []string `json:"tag"`
	UpdatedAt string   `json:"updated_at"`
}

func (h *NewsHandler) UpdateNews(c *fiber.Ctx) error {
	id := c.Params("id")
	var news UpdateNewsRequest

	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Update News input"})
	}

	fmt.Println(news.Category)

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ObjectID",
		})
	}

	existingNews, err := h.service.GetNewsByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve the news. Please try again later.",
		})
	}

	if news.Title != "" {
		existingNews.Title = news.Title
	}

	if news.Detail != "" {
		existingNews.Detail = news.Detail
	}

	if news.Image != "" {
		existingNews.Image = news.Image
	}

	if news.Category != "" {
		category, err := h.CategoryService.GetByID(news.Category)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve the news. Please try again later.",
			})
		}
		existingNews.CategoryID = &domain.Category{ID: category.ID, Name: category.Name}
	}

	if len(news.Tag) > 0 {
		existingNews.Tag = news.Tag
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	news.UpdatedAt = time.Now().In(loc).Format(time.RFC3339)

	err = h.service.UpdateNews(id, existingNews)
	if err != nil {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot Delete News cause Internal Server Error ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "News Delete Successfully",
	})
}
