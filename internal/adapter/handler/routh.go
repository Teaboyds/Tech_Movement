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
	Config             *config.HTTP
	NewsHandler        NewsHandler
	CategoryHandler    CategoryHandler
	MediaHandler       MediaHandler
	UploadHandler      UploadHandler
	InfographicHandler InfographicHandler
	BannerHandler      BannerHandler
}

func SetUpRoutes(p RouterParams) (*Router, error) {

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Static("/upload/news", "../upload/news")
	app.Static("/upload/media", "./upload/media")
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
			news.Get("/", p.NewsHandler.Find)
			news.Get("/:id", p.NewsHandler.GetNewsByID)
			// news.Get("/byCategory/:id", p.NewsHandler.GetNewsByCategory)
			news.Put("/:id", p.NewsHandler.UpdateNews)
			news.Delete("/:id", p.NewsHandler.DeleteNews)
			news.Delete("/", p.NewsHandler.DeleteManyNews)
		}

		category := v1.Group("/category")
		{
			category.Post("/", p.CategoryHandler.CreateCategory)
			category.Get("/:id", p.CategoryHandler.GetCategoryByID)
			category.Get("/", p.CategoryHandler.GetAllCategory)
			category.Put("/:id", p.CategoryHandler.UpdateCategory)
			category.Delete("/:id", p.CategoryHandler.DeleteCategory)
		}

		upload := v1.Group("/upload")
		{
			upload.Post("/", p.UploadHandler.UploadFile)
			upload.Get("/", p.UploadHandler.GetAllFile)
			upload.Get("/:id", p.UploadHandler.GetByID)
			upload.Delete("/:id", p.UploadHandler.DeleteFile)
		}

		media := v1.Group("/media")
		{
			media.Post("/", p.MediaHandler.CreateMedia)
		}

		infographic := v1.Group("/infographic")
		{
			infographic.Post("/", p.InfographicHandler.CreateInfographic)
			infographic.Get("/", p.InfographicHandler.GetInfoHome)
		}

		home := v1.Group("/home")
		{
			home.Get("/lastedNews", p.NewsHandler.GetLastNews)
			home.Get("/TechNews", p.NewsHandler.GetTechNews)
			home.Get("/", p.NewsHandler.GetHomePage)
			home.Get("/VDO", p.NewsHandler.GetVideoHome)
			home.Get("/Short", p.NewsHandler.GetShortVideoHome)
			home.Get("/Info", p.NewsHandler.GetInfoHome)
			// home.Get("/techNews", p.NewsHandler.GetNewsByCategory)
		}

		banner := v1.Group("/banner")
		{
			banner.Post("/", p.BannerHandler.CreateBanner)
			banner.Get("/:id", p.BannerHandler.GetBanner)
			banner.Get("/", p.BannerHandler.GetBanners)
			banner.Put("/:id", p.BannerHandler.UpdateBanner)
			banner.Delete("/:id", p.BannerHandler.DeleteBanner)
		}

	}
	return &Router{app}, nil
}

// func read address from config//
func (r *Router) Serve(listAddr string) error {
	return r.Listen(listAddr)
}
