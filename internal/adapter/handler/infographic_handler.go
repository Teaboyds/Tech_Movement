package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

type InfographicHandler struct {
	InfographicService port.InfographicService
}

func NewInfographicHandler(InfographicService port.InfographicService) *InfographicHandler {
	return &InfographicHandler{InfographicService: InfographicService}
}

func (ip *InfographicHandler) CreateInfographic(c *fiber.Ctx) error {

	info := new(domain.InfographicRequest)

	if err := c.BodyParser(info); err != nil {
		log.Printf("err: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Info input"})
	}

	if err := utils.FixInfoMultipart(c, info); err != nil {
		log.Printf("fix multipart array error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Form Array Field",
		})
	}

	if err := ip.InfographicService.CreateInfo(info); err != nil {
		log.Printf("errs: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Internal Server error : Create Info"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Infographic Create Successfully",
		"data":    info,
	})
}

func (ip *InfographicHandler) GetInfoHome(c *fiber.Ctx) error {
	info, err := ip.InfographicService.GetInfoHome()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth category data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "all categories",
		"data":    info,
	})
}
