package handler

import (
	dm "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	MediaService    port.MediaService
	CategoryService port.CategoryService
}

func NewMediaHandler(MediaService port.MediaService, CategoryService port.CategoryService) *MediaHandler {
	return &MediaHandler{
		MediaService:    MediaService,
		CategoryService: CategoryService,
	}
}

func (mh *MediaHandler) CreateMedia(c *fiber.Ctx) error {

	media := new(dm.MediaRequest)

	if err := c.BodyParser(media); err != nil {
		fmt.Printf("Bad request Create Media err: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error Bad Request Create Media",
		})
	}

	fileName, err := utils.UploadFile(c, "media_image", 5*1024*1024, "./upload/media_image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	tagList := strings.Split(media.Tag, ",")
	category, _ := mh.CategoryService.GetByID(media.Category)

	newMedia := dm.Media{
		Title:     media.Title,
		Abstract:  media.Abstract,
		Image:     fileName,
		Url:       media.Url,
		Category:  category,
		Tag:       tagList,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := mh.MediaService.CreateMedia(&newMedia); err != nil {
		fmt.Printf("Internal Create Media err: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Internal Server Error Create Media",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newMedia)
}

func (mh *MediaHandler) GetMediaHome(c *fiber.Ctx) error {

	lastMedia, err := mh.MediaService.GetLastMedia()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth lastMedia data",
		})
	}

	var lastedMedia []dm.HomeMediaResponse
	for _, media := range lastMedia {
		lastedMedia = append(lastedMedia, dm.HomeMediaResponse{
			Title:     media.Title,
			Abstract:  media.Abstract,
			Image:     media.Image,
			Url:       media.Url,
			Category:  media.Category.Name,
			CreatedAt: media.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "วีดีโอ",
		"data":    lastedMedia,
	})
}
