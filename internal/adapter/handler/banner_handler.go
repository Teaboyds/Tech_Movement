package handler

import (
	"backend_tech_movement_hex/internal/adapter/config"
	"backend_tech_movement_hex/internal/adapter/mapper"
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
	config        config.Container
}

func NewBannerHandler(bannerService port.BannerService, cateService port.CategoryService, config config.Container) *BannerHandler {
	return &BannerHandler{bannerService: bannerService, cateService: cateService, config: config}
}

func (ban *BannerHandler) CreateBanner(c *fiber.Ctx) error {

	var banner domain.BannerRequest

	banner.DesktopImage.FileType = "banner/desktop"
	banner.MobileImage.FileType = "banner/mobile"

	DesktopPath := banner.DesktopImage.FileType
	MobilePath := banner.MobileImage.FileType
	DesktopUploadDir := "./../upload/" + DesktopPath
	MobileUploadDir := "./../upload/" + MobilePath

	DesktopResult, err := utils.UploadFile(c, "desktop_image", 5*1024*1024, DesktopUploadDir)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	MobileResult, err := utils.UploadFile(c, "desktop_image", 5*1024*1024, MobileUploadDir)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	banner.DesktopImage.Name, banner.MobileImage.Name = DesktopResult.OriginalName, MobileResult.OriginalName
	banner.DesktopImage.Path, banner.MobileImage.Path = DesktopResult.SavedName, MobileResult.SavedName
	banner.DesktopImage.Type, banner.MobileImage.Type = DesktopResult.Type, MobileResult.Type

	if err := c.BodyParser(&banner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Errorf("banner bad request: %v", err),
		})
	}

	request := &domain.Banner{
		DesktopImage: banner.DesktopImage,
		MobileImage:  banner.MobileImage,
		Status:       banner.Status,
		LinkUrl:      banner.LinkUrl,
		Action:       banner.Action,
	}

	err = ban.bannerService.CreateBanner(request)
	if err != nil {
		log.Printf("create banner error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (banner *BannerHandler) GetBanner(c *fiber.Ctx) error {

	id := c.Params("id")

	bannerResp, err := banner.bannerService.GetBanner(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	bannerHandlerDTO := &domain.BannerResponse{
		ID:              bannerResp.ID,
		DesktopImageUrl: bannerResp.DesktopImage.Path,
		MobileImageUrl:  bannerResp.MobileImage.Path,
		Status:          bannerResp.Status,
		Action:          bannerResp.Action,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "banner retrive",
		"data":    bannerHandlerDTO,
	})
}

func (banner *BannerHandler) GetBanners(c *fiber.Ctx) error {
	query := c.Queries()
	visible := query["visible_page"]

	banners, err := banner.bannerService.GetBanners(visible)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var bannersResp []*domain.BannerResponse
	for _, banner := range banners {

		bannersHandlerDTO := &domain.BannerResponse{
			ID:              banner.ID,
			DesktopImageUrl: banner.DesktopImage.Path,
			MobileImageUrl:  banner.MobileImage.Path,
			Status:          banner.Status,
			Action:          banner.Action,
		}

		bannersResp = append(bannersResp, bannersHandlerDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "banner retrive",
		"data":    bannersResp,
	})
}

func (ban *BannerHandler) CreateBannerV2(c *fiber.Ctx) error {

	var banner domain.BannerRequestV2

	if err := c.BodyParser(&banner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	DirectoryDesktop := ban.config.IMG.BannerDesktop
	DirectoryMobile := ban.config.IMG.BannerMobile

	DesktopResult, err := utils.UploadFile(c, "desktop_image", 5*1024*1024, DirectoryDesktop)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	MobileResult, err := utils.UploadFile(c, "mobile_image", 5*1024*1024, DirectoryMobile)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	banner.DesktopImage = DesktopResult.SavedName
	banner.MobileImage = MobileResult.SavedName

	requestDTO := mapper.BannerRequestToDomain(banner, banner.Status)

	err = ban.bannerService.CreateBannerV2(requestDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(requestDTO)
}

func (ban *BannerHandler) GetBannerV2(c *fiber.Ctx) error {
	id := c.Params("id")

	bannerResp, err := ban.bannerService.GetBannerV2(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	DirectoryDesktop := ban.config.IMG.BannerDesktop
	DirectoryMobile := ban.config.IMG.BannerMobile
	DPath := ban.config.ImgPath.BannerDesktop
	MPath := ban.config.ImgPath.BannerMobile

	parseDestop, err := utils.GetImageMetadata(DirectoryDesktop, bannerResp.DesktopImage, DPath)
	if err != nil {
		return nil
	}

	parseMobile, err := utils.GetImageMetadata(DirectoryMobile, bannerResp.MobileImage, MPath)
	if err != nil {
		return nil
	}

	BannerResponse := mapper.BannerDomainToResponse(*bannerResp, *parseDestop, *parseMobile, bannerResp.Status)

	fmt.Printf("BannerResponse: %v\n", BannerResponse)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "banner retrive",
		"data":    BannerResponse,
	})
}
