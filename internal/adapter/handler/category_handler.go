package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryHandler struct {
	CategoryService port.CategoryService
}

func NewCategoryHandler(CategoryService port.CategoryService) *CategoryHandler {
	return &CategoryHandler{CategoryService: CategoryService}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	category := new(domain.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Category input"})
	}

	if err := h.CategoryService.Create(category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category ",
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

	// ใช้วรยุทธ์สร้าง struct ร่างโคลนอีกหนึ่งร่าง เพื่อรับ body request เฉพาะฝั่ง http โดยเฉพาะ เพื่อให้แยกออกจากฝั่ง domain โดยเด็ดขาด //
	type CategoryResponse struct {
		ID   primitive.ObjectID `json:"id"`
		Name string             `json:"name"`
	}

	// instance categories ให้เป็น slice CategoryResponse แล้ว append เข้าไป //
	var categories []CategoryResponse
	for _, catcategories := range category {
		categories = append(categories, CategoryResponse{
			ID:   catcategories.ID,
			Name: catcategories.Name,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "all categories",
		"data":    categories,
	})
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	category := new(domain.Category)

	if err := c.BodyParser(category); err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Category Input",
		})
	}

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Category Input",
		})
	}

	// วรยุทธเดียวกันกับตรง get //
	NewCategory := &domain.Category{
		ID:   ObjID,
		Name: category.Name,
	}

	err = h.CategoryService.UpdateCategory(ObjID.Hex(), NewCategory)
	if err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category Updated",
		"data":    NewCategory,
	})
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Category ID Input",
		})
	}

	err = h.CategoryService.DeleteCategory(ObjID.Hex())
	if err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot Delete Category cause Database Issue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category Delete Successfully",
	})
}
