package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	medias "backend-tech-movement/models/Medias"
	news "backend-tech-movement/models/News"
	banners "backend-tech-movement/models/banners"
	infographics "backend-tech-movement/models/infographics"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func LoadEnv(){
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Env coudln't Load",err)
	}
}

func ConnectDB() *gorm.DB { 

	LoadEnv()

	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	fmt.Println(portStr)

	port , err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Invalid DB port", err)
	}

	dns := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: newLogger,
	})
	
	db.AutoMigrate(&news.NewsCategory{}, &news.NewsTag{}, &news.News{} , &news.NewsDetail{})
	db.AutoMigrate(&medias.MediaCategory{}, &medias.MediaTag{}, &medias.Medias{})
	db.AutoMigrate(&banners.BannersCategory{}, &banners.Banners{})
	db.AutoMigrate(&infographics.Infographics{}, &infographics.InfographicDetails{})


	if err != nil {
		panic("failed to connect to database")
	}

	return db
}