package newscontrollers

import (
	news "backend-tech-movement/models/News"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// create news //
func CreateNews(db *gorm.DB) fiber.Handler {
	
	return func(c *fiber.Ctx) error {

	newsData := &news.News{}
	if err := c.BodyParser(newsData); err != nil{
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc).Format(time.RFC3339)
	newsData.CreatedAt = now
	newsData.UpdatedAt = now

	if err := db.Create(newsData).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})	
	}

	return c.JSON(fiber.Map{
		"status":fiber.StatusCreated ,
		"data": newsData,
	})
	}
}

// update news //
func UpdateNews(db *gorm.DB) fiber.Handler {
	
	return func(c *fiber.Ctx) error {
		// find data //
	id := c.Params("id")
	newsData := &news.News{}
	if err := db.First(newsData , "id = ?" , id); err != nil{
		log.Println(err)
		return c.Status(404).JSON(fiber.Map{
			"message": "News not Found.",
		})	
	}

	// request body kubpom //
	inputData := &news.News{}
	if err := c.BodyParser(inputData);err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})	
	}

	// set time eiei // 
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc).Format(time.RFC3339)
	inputData.UpdatedAt = now

	//  db save and update //
	if err := db.Model(&newsData).Updates(inputData).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Internal Server Error",
		})	
	}

	return c.JSON(fiber.Map{
		"status":fiber.StatusOK ,
		"data": newsData,
	})
	}
}

func DeleteNews(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		newsData := &news.News{}

		// news checker//
		if err := db.First(newsData , "id = ?" , id).Error; err != nil{
			log.Println(err)
			return c.Status(404).JSON(fiber.Map{
				"message": "News not Found.",
			})	
		}

		// delete func //
		if err := db.Delete(newsData , "id = ?",id).Error; err != nil {
			log.Println(err)
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})	
		}

		return c.JSON(fiber.Map{
			"status":fiber.StatusOK ,
		})
	}
}

// create detail one to one on news // 
func CreateDetail(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		detail := &news.NewsDetail{}

		if err := c.BodyParser(&detail); err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Bad Request",
			})			
		}

		if err := db.Create(&detail).Error; err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})				
		}

			return c.JSON(fiber.Map{
				"status":fiber.StatusCreated ,
				"data": detail,
			})
	}
}