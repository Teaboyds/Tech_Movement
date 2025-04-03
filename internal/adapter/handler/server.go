// routh //

package handler

import (
	"backend_tech_movement_hex/internal/adapter/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

type Router struct {
	*fiber.App
}

type RouterParams struct {
	Config          *config.HTTP
	NewsHandler     NewsHandler
	CategoryHandler CategoryHandler
}

func SetUpRoutes(p RouterParams) (*Router, error) {

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     p.Config.HttpOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	api := app.Group(fmt.Sprintf("%s%s", p.Config.Prefix, "/api"))

	v1 := api.Group("/v1")
	{
		news := v1.Group("/news")
		{
			news.Post("/", p.NewsHandler.CreateNews)
			news.Get("/:id", p.NewsHandler.GetNewsByID)
			news.Get("/category/:id", p.NewsHandler.GetNewsByCategory)
			news.Get("/", p.NewsHandler.GetNewsByPage)
			news.Put("/:id", p.NewsHandler.UpdateNews)
			news.Delete("/:id", p.NewsHandler.DeleteNews)
			news.Get("/uploads/:filename", func(c *fiber.Ctx) error {
				filename := c.Params("filename")
				directory := "./internal/core/upload"
				filePath := directory + "/" + filename

				return c.SendFile(filePath)
			})
		}

		category := v1.Group("/category")
		{
			category.Post("/", p.CategoryHandler.CreateCategory)
			category.Get("/:id", p.CategoryHandler.GetCategoryByID)
			category.Get("/", p.CategoryHandler.GetAllCategory)
			category.Put("/:id", p.CategoryHandler.UpdateCategory)
			category.Delete("/:id", p.CategoryHandler.DeleteCategory)
		}
	}

	return &Router{app}, nil
}

// func read address from config//
func (r *Router) Serve(listAddr string) error {
	return r.Listen(listAddr)
}
