package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type BannerHandler struct {
	bannerService port.BannerService
	cateService   port.CategoryService
}

func NewBannerHandler(bannerService port.BannerService, cateService port.CategoryService) *BannerHandler {
	return &BannerHandler{bannerService: bannerService, cateService: cateService}
}

func (ban *BannerHandler) CreateBanner(c *fiber.Ctx) error {

	var banner domain.BannerReq

	if err := c.BodyParser(&banner); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fieldMap := map[string]string{
		"image": "Img",
	}

	if err := utils.FixMultipartArrayV2(c, &banner, fieldMap); err != nil {
		log.Printf("fix multipart array error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Form Array Field",
		})
	}

	jeans := &domain.Banner{
		Title:       banner.Title,
		ContentType: banner.ContentType,
		Status:      banner.Status,
		Category:    banner.Category,
		Img:         banner.Img,
	}

	fmt.Printf("jeans: %v\n", jeans)

	if err := ban.bannerService.CreateBanner(jeans); err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Banner Create Successfully",
		"data":    jeans,
	})
}

func (ban *BannerHandler) GetBanner(c *fiber.Ctx) error {

	id := c.Params("id")

	banner, err := ban.bannerService.GetBanner(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrResponse{
			Error: "Category Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category",
		"data":    banner,
	})

}

func (ban *BannerHandler) GetBanners(c *fiber.Ctx) error {
	banners, err := ban.bannerService.GetBanners()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrResponse{
			Error: "Category Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category",
		"data":    banners,
	})

}

func (ban *BannerHandler) UpdateBanner(c *fiber.Ctx) error {
	id := c.Params("id")

	var banner domain.BannerReq
	if err := c.BodyParser(&banner); err != nil {
		log.Printf("banner bad request %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Banner Reuquest",
		})
	}
	fieldMap := map[string]string{
		"image": "Img",
	}

	if err := utils.FixMultipartArrayV2(c, &banner, fieldMap); err != nil {
		log.Printf("fix multipart array error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Form Array Field",
		})
	}
	fmt.Printf("banner.Img: %v\n", banner.Img)
	bannernew := &domain.Banner{
		Title:       banner.Title,
		ContentType: banner.ContentType,
		Status:      banner.Status,
		Category:    banner.Category,
		Img:         banner.Img,
	}
	if err := ban.bannerService.Updated(id, bannernew); err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot Update Banner cause Internal Server Error ",
		})
	}
	fmt.Printf("bannernew.Img: %v\n", bannernew.Img)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Banner Updated Successfully",
		"data":    bannernew,
	})

}

func (ban *BannerHandler) DeleteBanner(c *fiber.Ctx) error {
	id := c.Params("id")

	err := ban.bannerService.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrResponse{
			Error: "Banner Not Found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Banner Deleted Successfully",
	})
}
