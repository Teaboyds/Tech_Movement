package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func LoadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Env coudln't Load", err)
	}
}

func ConnectedRedis() {

	LoadEnv()

	RedisPort := os.Getenv("REDIS_PORT")
	if RedisPort == "" {
		RedisPort = "localhost:6379"
	}
	RedisPass := os.Getenv("REDIS_PASSWORD")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisPort,
		Password: RedisPass,
		DB:       0,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Cannot connect to Redis: ", err)
	}
	log.Println("Connected to Redis")
}

func GetRedisClient() *redis.Client {
	return RedisClient
}
