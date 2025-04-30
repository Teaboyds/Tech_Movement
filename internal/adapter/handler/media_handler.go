package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"

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

	if err := utils.FixMediaMultipart(c, media); err != nil {
		log.Printf("fix multipart array error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Form Array Field",
		})
	}

	fmt.Printf("media: %v\n", media)

	if err := m.MediaService.CreateMedia(media); err != nil {
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
