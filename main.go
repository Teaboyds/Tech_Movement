package main

import (
	_ "backend_tech_movement_hex/docs"
	"backend_tech_movement_hex/internal/adapter/handler"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/repository"
	"backend_tech_movement_hex/internal/adapter/storage/redis"
	routh "backend_tech_movement_hex/internal/app"
	"backend_tech_movement_hex/internal/core/service"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title News API
// @description This is a sample server for a News API.
// @version 1.0
// @host localhost:5050
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000 , http://127.0.0.1:3000",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)
	redis.ConnectedRedis()

	// debug : if redis disconnected //
	if redis.RedisClient == nil {
		log.Fatal("❌ RedisClient is still nil after ConnectedRedis()")
	} else {
		log.Println("✅ RedisClient initialized:", redis.RedisClient)
	}

	mongodb.ConnectDB()

	//Redis//
	cacheRepo := redis.NewRedisCacheRepository(redis.RedisClient)

	// news //
	categoryRepo := repository.NewCategoryRepositoryMongo(mongodb.GetDatabase())
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	newsRepo := repository.NewNewsRepo(mongodb.GetDatabase(), redis.RedisClient)
	newService := service.NewsService(newsRepo, categoryRepo)
	newHandler := handler.NewNewsHandler(newService, categoryRepo, cacheRepo)

	routh.SetUpRoutes(app, newHandler, categoryHandler)
	app.Listen(":5050")
}
