// routh //

package routh

import (
	"backend_tech_movement_hex/internal/adapter/handler"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App, newsHandler *handler.NewsHandler) {

	news := app.Group("/news")
	{
		news.Post("/", newsHandler.CreateNews)
		news.Get("/:id", newsHandler.GetNewsByID)
		news.Get("/", newsHandler.GetNewsByPage)
		news.Put("/:id", newsHandler.UpdateNews)
		news.Delete("/:id", newsHandler.DeleteNews)
		news.Get("/uploads/:filename", func(c *fiber.Ctx) error {
			filename := c.Params("filename")
			directory := "./internal/core/upload"
			filePath := directory + "/" + filename

			return c.SendFile(filePath)
		})
	}

}
