package newscontrollers

import (
	news "backend-tech-movement/models/News"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateCategory(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		category := &news.NewsCategory{}
		if err := c.BodyParser(&category); err != nil{
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Bad Request",
			})			
		}
		

		loc, _ := time.LoadLocation("Asia/Bangkok")
		now := time.Now().In(loc).Format(time.RFC3339)
		category.CreatedAt = now
		category.UpdatedAt = now

		if err := db.Create(&category).Error; err != nil{
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})	
		}

		return c.JSON(fiber.Map{
			"status":fiber.StatusCreated ,
			"data": category,
		})
	}
}