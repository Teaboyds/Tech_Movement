package main

import (
	_ "backend_tech_movement_hex/docs"
	"backend_tech_movement_hex/internal/adapter/config"
	"backend_tech_movement_hex/internal/adapter/handler"
	"backend_tech_movement_hex/internal/adapter/logger"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/repository"
	cache "backend_tech_movement_hex/internal/adapter/storage/redis"
	"backend_tech_movement_hex/internal/core/service"
	"context"
	"fmt"
	"log/slog"
	"os"
)

func Init(config *config.Container) {
	logger.ZapInit()

	ctx := context.Background()

	fmt.Println("MongoDB URI:", config.DB.URL)
	fmt.Println("MongoDB Name:", config.DB.DB_NAME)

	// connect mongodb //
	db, err := mongodb.ConnectDB(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	// connect redis //
	// redisOptions := &redis.Options{
	// 	Addr:     config.Redis.REDIS_PORT,
	// 	Password: config.Redis.REDIS_PASSWORD,
	// 	DB:       config.Redis.REDIS_DB,
	// }
	// redisClient := redis.NewClient(redisOptions)
	cacheClient, err := cache.ConnectedRedis(ctx, config.Redis)
	if err != nil {
		slog.Error("Error initializing cache connection", "error", err)
		os.Exit(1)
	}

	// Category //
	categoryRepo := repository.NewCategoryRepositoryMongo(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// News //
	newsRepo := repository.NewNewsRepo(db)
	if err := newsRepo.EnsureNewsIndexs(); err != nil {
		slog.Error("Error ensuring news indexes", "error", err)
		os.Exit(1)
	}
	newService := service.NewsService(newsRepo, categoryRepo, cacheClient)
	newHandler := handler.NewNewsHandler(newService, categoryRepo, cacheClient)

	mediaRepo := repository.NewMediaRepo(db)
	if err := mediaRepo.EnsureMediaIndexs(); err != nil {
		slog.Error("Error ensuring news indexes", "error", err)
		os.Exit(1)
	}
	mediaService := service.NewMediaService(mediaRepo)
	mediaHandler := handler.NewMediaHandler(mediaService, categoryService)

	// run server from server.go //
	router, err := handler.SetUpRoutes(handler.RouterParams{
		Config:          config.HTTP,
		NewsHandler:     *newHandler,
		CategoryHandler: *categoryHandler,
		MediaHandler:    *mediaHandler,
	})
	if err != nil {
		logger.Error("Error initializing router" + "error" + err.Error())
		os.Exit(1)
	}

	//read address from config//
	logger.Info(config.HTTP.URL + ":" + config.HTTP.Port)
	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	logger.Info("Starting the HTTP server " + "< listen_address : " + listenAddr + " >")
	err = router.Serve(listenAddr)
	if err != nil {
		logger.Error("Error starting the HTTP server" + "error" + err.Error())
		os.Exit(1)
	}

}

// @description This is a sample server for a News API.
// @version 1.0
// @host localhost:5050
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	config, err := config.New()
	if err != nil {
		return
	}

	Init(config)

}
