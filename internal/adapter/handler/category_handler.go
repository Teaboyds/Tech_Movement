package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	CategoryService port.CategoryService
}

func NewCategoryHandler(CategoryService port.CategoryService) *CategoryHandler {
	return &CategoryHandler{CategoryService: CategoryService}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {

	category := new(domain.CategoryRequest)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Category input"})
	}

	if err := utils.ValidateCategoryInput(category); err != nil {
		log.Printf("news bad validator request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrResponse{
			Error: "Invalid Validator News Request",
		})
	}

	cateDTO := &domain.Category{
		Name:         category.Name,
		CategoryType: category.CategoryType,
	}

	if err := h.CategoryService.Create(cateDTO); err != nil {
		fmt.Printf("err: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category Create Successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {

	id := c.Params("id")

	category, err := h.CategoryService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ErrResponse{
			Error: "Category Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category",
		"data":    category,
	})

}

func (h *CategoryHandler) GetAllCategory(c *fiber.Ctx) error {

	category, err := h.CategoryService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth category data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "all categories",
		"data":    category,
	})
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category domain.Category

	if err := c.BodyParser(&category); err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Category Input",
		})
	}

	if err := h.CategoryService.UpdateCategory(id, &category); err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category Updated",
		"data":    category,
	})
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.CategoryService.DeleteCategory(id); err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot Delete Category cause Database Issue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category Delete Successfully",
	})
}
