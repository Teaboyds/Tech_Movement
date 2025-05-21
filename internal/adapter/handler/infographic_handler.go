package handler

import (
	"backend_tech_movement_hex/internal/adapter/mapper"
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

type InfographicHandler struct {
	InfographicService port.InfographicService
	CategoryService    port.CategoryService
	ImageService       port.UploadService
}

func NewInfographicHandler(InfographicService port.InfographicService, CategoryService port.CategoryService, ImageService port.UploadService) *InfographicHandler {
	return &InfographicHandler{InfographicService: InfographicService, CategoryService: CategoryService, ImageService: ImageService}
}

func (ip *InfographicHandler) CreateInfographic(c *fiber.Ctx) error {

	info := new(domain.InfographicRequestDTO)

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

	infoRequest := mapper.InfographicRequestToDomain(info)

	if err := ip.InfographicService.CreateInfo(infoRequest); err != nil {
		log.Printf("errs: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Internal Server error : Create Info"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Infographic Create Successfully",
		"data":    info,
	})
}

func (ip *InfographicHandler) GetInfographic(c *fiber.Ctx) error {

	id := c.Params("id")

	infographic, err := ip.InfographicService.GetInfographic(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error()},
		)
	}

	caegory, err := ip.CategoryService.GetByID(infographic.Category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "cannot fetch category in get infographic",
		})
	}

	image, err := ip.ImageService.GetFileByID(infographic.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "cannot fetch image in get infographic",
		})
	}

	infographicResp := mapper.EnrichInfographic(infographic, caegory, image)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Infographic Create Successfully",
		"data":    infographicResp,
	})
}

func (ip *InfographicHandler) GetInfographics(c *fiber.Ctx) error {

	query := c.Queries()
	cateId := query["category"]
	sort := query["sort"]
	view := query["view"]
	limit := query["limit"]
	page := query["page"]

	infographic, err := ip.InfographicService.GetInfographics(cateId, sort, view, limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "cannot fecth infographic",
		})
	}

	categoryImageIds := mapper.ExtractCategoryAndImageIDs(infographic)

	category, err := ip.CategoryService.GetByIDs(categoryImageIds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	image, err := ip.ImageService.GetFilesByIDsVTest(categoryImageIds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	CategoryMap := make(map[string]domain.Category)
	for _, cat := range category {
		CategoryMap[cat.ID] = *cat
	}

	ImgMap := make(map[string]domain.UploadFile)
	for _, img := range image {
		ImgMap[img.ID] = *img
	}

	infoResp := mapper.EnrichInfographics(infographic, CategoryMap, ImgMap)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "media retrive successfully",
		"data":    infoResp,
	})
}

// func (ip *InfographicHandler) GetInfoHome(c *fiber.Ctx) error {
// 	info, err := ip.InfographicService.GetInfoHome()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "cannot fecth category data",
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "all categories",
// 		"data":    info,
// 	})
// }
