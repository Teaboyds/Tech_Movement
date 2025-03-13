package routes

import (
	"backend-tech-movement/controllers"
	newscontrollers "backend-tech-movement/controllers/news_controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(c *fiber.App) {
	api := c.Group("/api")

	api.Get("/users", controllers.GetUsers)
	api.Post("/users", controllers.CreateUser)
}

func SetupNewsRoutes(c *fiber.App  , db *gorm.DB){
	api := c.Group("/news")

	api.Post("/", newscontrollers.CreateNews(db))
	api.Put("/",  newscontrollers.UpdateNews(db))
}

func SetupCategoryRoutes(c *fiber.App  , db *gorm.DB){
	api := c.Group("/category")

	api.Post("/" , newscontrollers.CreateCategory(db) )
}