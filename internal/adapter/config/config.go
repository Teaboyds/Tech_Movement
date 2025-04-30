package config

import (
	"backend_tech_movement_hex/internal/core/utils"
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App   *App
		DB    *DB
		HTTP  *HTTP
		Redis *Redis
		OneID *OneID
		JWT   *JWT
	}

	App struct {
		Name       string
		Production bool
		Version    string
	}

	DB struct {
		URL     string
		DB_NAME string
	}

	HTTP struct {
		Env            string
		URL            string
		Port           string
		HttpOrigins    string
		AllowedOrigins string
		Prefix         string
	}

	Redis struct {
		REDIS_HOST     string
		REDIS_PORT     string
		REDIS_PASSWORD string
		REDIS_DB       int
	}

	OneID struct {
		URL          string
		ClientID     string
		ClientSecret string
		RedirectURL  string
	}

	JWT struct {
		SecretKey string
	}
)

func New() (*Container, error) {
	if utils.GetEnv("APP_ENV", "development") != "production" {
		if _, err := os.Stat("./../.env"); err == nil {
			err := godotenv.Load("./../.env")
			if err != nil {
				return nil, err
			}
		}
	}

	app := &App{
		Name:       utils.GetEnv("APP_NAME", ""),
		Production: utils.GetEnv("APP_PRODUCTION", "false") == "true",
		Version:    utils.GetEnv("APP_VERSION", "0.0.0"),
	}

	db := &DB{
		URL:     utils.GetEnv("MONGO_URL", ""),
		DB_NAME: utils.GetEnv("DB_NAME_MONGOD", ""),
	}

	http := &HTTP{
		Env:            utils.GetEnv("APP_ENV", ""),
		URL:            utils.GetEnv("HTTP_URL", ""),
		Port:           utils.GetEnv("HTTP_PORT", ""),
		HttpOrigins:    utils.GetEnv("HTTP_ORIGINS", ""),
		AllowedOrigins: utils.GetEnv("HTTP_ALLOWED_ORIGINS", ""),
	}

	redis := &Redis{
		REDIS_HOST:     utils.GetEnv("REDIS_HOST", ""),
		REDIS_PORT:     utils.GetEnv("REDIS_PORT", ""),
		REDIS_PASSWORD: utils.GetEnv("REDIS_PASSWORD", ""),
		REDIS_DB:       utils.AtoI(utils.GetEnv("REDIS_DB", "0"), 0),
	}

	oneID := &OneID{
		URL:          utils.GetEnv("ONE_URL_UAT", ""),
		ClientID:     utils.GetEnv("ONE_ID_CLIENT_ID_UAT", ""),
		ClientSecret: utils.GetEnv("ONE_ID_CLIENT_SECRET_UAT", ""),
		RedirectURL:  utils.GetEnv("ONE_ID_REDIRECT_URL_UAT", ""),
	}

	jwt := &JWT{
		SecretKey: utils.GetEnv("JWT_SECRET_KEY", ""),
	}

	return &Container{
		App:   app,
		DB:    db,
		HTTP:  http,
		Redis: redis,
		OneID: oneID,
		JWT:   jwt,
	}, nil

}
