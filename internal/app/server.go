// routh //

package routh

import (
	"backend_tech_movement_hex/internal/adapter/handler"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App, newsHandler *handler.NewsHandler, categoryHandler *handler.CategoryHandler) {

	news := app.Group("/news")
	{
		news.Post("/", newsHandler.CreateNews)
		news.Get("/:id", newsHandler.GetNewsByID)
		news.Get("/", newsHandler.GetNewsByPage)
		news.Put("/:id", newsHandler.UpdateNews)
		news.Delete("/:id", newsHandler.DeleteNews)
		news.Options("/*", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNoContent)
		})
		news.Get("/uploads/:filename", func(c *fiber.Ctx) error {
			filename := c.Params("filename")
			directory := "./internal/core/upload"
			filePath := directory + "/" + filename

			return c.SendFile(filePath)
		})
	}

	category := app.Group("/category")
	{
		category.Post("/", categoryHandler.CreateCategory)
		category.Get("/:id", categoryHandler.GetCategoryByID)
		category.Get("/", categoryHandler.GetAllCategory)
		category.Put("/:id", categoryHandler.UpdateCategory)
		category.Delete("/:id", categoryHandler.DeleteCategory)
	}

}
