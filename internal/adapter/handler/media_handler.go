package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	MediaService port.MediaService
}

func NewMediaHandler(MediaService port.MediaService) *MediaHandler {
	return &MediaHandler{MediaService: MediaService}
}

func (m *MediaHandler) CreateMedia(c *fiber.Ctx) error {

	media := new(domain.MediaRequest)

	if err := c.BodyParser(media); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Media input"})
	}

	req := &domain.Media{
		Title:      media.Title,
		Content:    media.Content,
		VideoURL:   media.VideoURL,
		ThumnailID: media.ThumnailID,
		CategoryID: media.CategoryID,
		Tags:       media.Tags,
		Action:     media.Action,
	}

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
