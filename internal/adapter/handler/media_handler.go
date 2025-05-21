package handler

import (
	"backend_tech_movement_hex/internal/adapter/mapper"
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	MediaService    port.MediaService
	CategoryService port.CategoryService
}

func NewMediaHandler(MediaService port.MediaService, CategoryService port.CategoryService) *MediaHandler {
	return &MediaHandler{MediaService: MediaService, CategoryService: CategoryService}
}

func (m *MediaHandler) CreateMedia(c *fiber.Ctx) error {

	media := new(domain.MediaRequest)

	if err := c.BodyParser(media); err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Media input"})
	}

	req := mapper.MediaRequestToDomain(*media)

	if err := m.MediaService.CreateMedia(req); err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create media ",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Media Create Successfully",
		"data":    media,
	})
}

func (m *MediaHandler) GetMedias(c *fiber.Ctx) error {

	query := c.Queries()
	cateId := query["category"]
	sort := query["sort"]
	view := query["view"]
	limit := query["limit"]
	page := query["page"]

	media, err := m.MediaService.GetMedias(cateId, sort, view, limit, page)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	categoryIDs := mapper.ExtractCategoryIDs(media)
	categories, err := m.CategoryService.GetByIDs(categoryIDs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	categoryMap := make(map[string]domain.Category)
	for _, cat := range categories {
		categoryMap[cat.ID] = *cat
	}

	mediasResp := mapper.EnrichMedias(media, categoryMap)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "media retrive successfully",
		"data":    mediasResp,
	})
}

func (m *MediaHandler) GetMedia(c *fiber.Ctx) error {
	id := c.Params("id")

	media, err := m.MediaService.GetMedia(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	category, err := m.CategoryService.GetByID(media.Category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	mediaResp := mapper.EnrichMedia(media, category)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "media retrive successfully",
		"data":    mediaResp,
	})
}
