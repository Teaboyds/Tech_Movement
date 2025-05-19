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

	// Upload //
	uploadRepo := repository.NewUploadRepo(db)
	if err := uploadRepo.EnsureFileIndexs(); err != nil {
		slog.Error("Error ensuring news indexes", "error", err)
		os.Exit(1)
	}
	uploadServiec := service.NewUploadService(uploadRepo)
	uploadHander := handler.NewUploadHandler(uploadServiec)

	// Media //
	mediaRepo := repository.NewMediaRepositoryMongo(db)
	mediaService := service.NewMediaService(mediaRepo, categoryRepo, categoryService, uploadServiec)
	mediaHandler := handler.NewMediaHandler(mediaService)

	// InfoGraphic //
	infographicRepo := repository.NewInfographicRepositoryMongo(db)
	infographicService := service.NewInfographicService(infographicRepo, uploadServiec, categoryService)
	infographicHandler := handler.NewInfographicHandler(infographicService)

	// Banner //
	bannerRepo := repository.NewBannersRepoMongo(db)
	bannerService := service.NewBannerService(bannerRepo, categoryRepo, uploadServiec, uploadRepo, categoryService)
	bannerHandler := handler.NewBannerHandler(bannerService, categoryService)

	// News //
	newsRepo := repository.NewNewsRepo(db)
	if err := newsRepo.EnsureNewsIndexs(); err != nil {
		slog.Error("Error ensuring news indexes", "error", err)
		os.Exit(1)
	}
	newService := service.NewsService(newsRepo, categoryRepo, cacheClient, uploadRepo, categoryService, uploadServiec)
	newHandler := handler.NewNewsHandler(mediaService, newService, categoryService, cacheClient, infographicService)

	// run server from server.go //
	router, err := handler.SetUpRoutes(handler.RouterParams{
		Config:             config.HTTP,
		NewsHandler:        *newHandler,
		CategoryHandler:    *categoryHandler,
		MediaHandler:       *mediaHandler,
		UploadHandler:      *uploadHander,
		InfographicHandler: *infographicHandler,
		BannerHandler:      *bannerHandler,
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
