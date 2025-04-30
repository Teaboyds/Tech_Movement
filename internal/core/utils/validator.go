package utils

import (
	"backend_tech_movement_hex/internal/core/domain"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateNewsInput(input *domain.NewsRequest) error {
	if err := validate.Struct(input); err != nil {
		return err
	}
	return nil
}

func ValidateUploadInput(file *domain.UploadFileRequest) error {
	if err := validate.Struct(file); err != nil {
		return err
	}
	return nil
}

func ValidateMediaInput(file *domain.MediaRequest) error {
	if err := validate.Struct(file); err != nil {
		return err
	}
	return nil
}

func ValidateCategoryInput(file *domain.CategoryRequest) error {
	if err := validate.Struct(file); err != nil {
		return err
	}
	return nil
}

func FixInfoMultipart(c *fiber.Ctx, input *domain.InfographicRequest) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	if catID, ok := form.Value["category_id"]; ok && len(catID) > 0 {
		input.Category = catID[0]
	}

	return nil
}

func FixMediaMultipart(c *fiber.Ctx, input *domain.MediaRequest) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	if catID, ok := form.Value["category_id"]; ok && len(catID) > 0 {
		input.Category = catID[0]
	}

	return nil
}

func FixMultipartArray(c *fiber.Ctx, input *domain.NewsRequest) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// ถ้า form มี field news_image หรือ tag ถึงค่อย override
	if newsImgs, ok := form.Value["news_image"]; ok {
		input.Image = newsImgs
	}

	return nil
}
