package handler

import (
	"backend_tech_movement_hex/internal/core/domain"
	dt "backend_tech_movement_hex/internal/core/domain"
	p "backend_tech_movement_hex/internal/core/port"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagsHandler struct {
	TagService p.TagsService
}

func NewTagsHandler(TagService p.TagsService) *TagsHandler {
	return &TagsHandler{TagService: TagService}
}

func (tr *TagsHandler) CreateTags(c *fiber.Ctx) error {

	tags := new(dt.Tags)

	if err := c.BodyParser(tags); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "Error Bad Request",
		})
	}

	if err := tr.TagService.CreateTags(tags); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tags Create Successfully",
		"data":    tags,
	})
}

func (tr *TagsHandler) GetTagsByID(c *fiber.Ctx) error {

	id := c.Params("id")

	tags, err := tr.TagService.GetByID(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Can't Convert ID",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tags",
		"data":    tags,
	})
}

func (tr *TagsHandler) GetAll(c *fiber.Ctx) error {

	tag, err := tr.TagService.GetTagsAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot fecth tags data",
		})
	}

	var tagsresponse []dt.TagsResponse
	for _, tags := range tag {
		tagsresponse = append(tagsresponse, dt.TagsResponse{
			ID:   tags.ID,
			Name: tags.Name,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "all tags",
		"data":    tag,
	})
}

func (tr *TagsHandler) UpdateTag(c *fiber.Ctx) error {
	id := c.Params("id")
	tag := new(domain.Tags)

	if err := c.BodyParser(tag); err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Tag Input",
		})
	}

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Tag Input",
		})
	}

	// วรยุทธเดียวกันกับตรง get //
	NewTag := &domain.Tags{
		ID:   ObjID,
		Name: tag.Name,
	}

	err = tr.TagService.UpdateTags(ObjID.Hex(), NewTag)
	if err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update Tag",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tag Updated",
		"data":    NewTag,
	})
}

func (tr *TagsHandler) DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")
	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Tag ID Input",
		})
	}

	err = tr.TagService.DeleteTags(ObjID.Hex())
	if err != nil {
		log.Printf("Error is %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot Delete Tag cause Database Issue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tag Delete Successfully",
	})
}
