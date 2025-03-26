package main

import (
	_ "backend_tech_movement_hex/docs"
	"backend_tech_movement_hex/internal/adapter/handler"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/repository"
	routh "backend_tech_movement_hex/internal/app"
	"backend_tech_movement_hex/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title News API
// @description This is a sample server for a News API.
// @version 1.0
// @host localhost:5050
// @BasePath /api/news
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	mongodb.ConnectDB()

	// news //
	categoryRepo := repository.NewCategoryRepositoryMongo(mongodb.GetDatabase())
	newsRepo := repository.NewNewsRepo(mongodb.GetDatabase())
	newService := service.NewsService(newsRepo, categoryRepo)
	newHandler := handler.NewNewsHandler(newService, categoryRepo)
	routh.SetUpRoutes(app, newHandler)
	app.Listen(":5050")
}
