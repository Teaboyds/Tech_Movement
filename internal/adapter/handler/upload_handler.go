package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	uploadService port.UploadService
}

func NewUploadHandler(uploadService port.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

func (ul *UploadHandler) UploadFile(c *fiber.Ctx) error {

	var file domain.UploadFileRequest
	if err := c.BodyParser(&file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file upload"})
	}

	fmt.Printf("file: %v\n", file)

	if err := utils.ValidateUploadInput(&file); err != nil {
		log.Printf("uploadFile bad validator request %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Validator Upload Reuquest",
		})
	}

	path := file.FileType
	uploadDir := "./../upload/" + path
	fmt.Printf("uploadDir: %v\n", uploadDir)

	result, err := utils.UploadFile(c, "file", 5*1024*1024, uploadDir)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	file.Name = result.OriginalName
	file.Path = result.SavedName

	if err := ul.uploadService.UploadFile(&file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "File Upload Successfully",
		"data":    file,
	})
}

func (ul *UploadHandler) GetAllFile(c *fiber.Ctx) error {

	files, err := ul.uploadService.GetAllFile()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth all files",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "all categories",
		"data":    files,
	})
}

func (ul *UploadHandler) GetByID(c *fiber.Ctx) error {

	fileID := c.Params("id")

	file, err := ul.uploadService.GetFileByID(fileID)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to find the requested file.",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": file.FileType,
		"data":    file,
	})
}

func (ul *UploadHandler) DeleteFile(c *fiber.Ctx) error {

	id := c.Params("id")

	if err := ul.uploadService.DeleteFile(id); err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot Delete File cause Database Issue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File Delete Successfully",
	})

}
