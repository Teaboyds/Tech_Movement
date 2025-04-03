package redis

import (
	"backend_tech_movement_hex/internal/adapter/config"
	"backend_tech_movement_hex/internal/core/port"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	MongoClient *redis.Client
}

func ConnectedRedis(ctx context.Context, config *config.Redis) (port.CacheRepository, error) {
	MongoClient := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_PORT,
		Password: config.REDIS_PASSWORD,
		DB:       config.REDIS_DB,
	})
	_, err := MongoClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Cannot connect to Redis: ", err)
	}
	log.Println("Connected to Redis")
	return &Redis{MongoClient}, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, timeout time.Duration) error {

	// แปลง value เพื่อเก็บใน redis เพราะ ฝ //
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err = r.MongoClient.Set(ctx, key, j, timeout).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {

	// ดึง keys มาแล้วแปลงเป็น byte //
	b, err := r.MongoClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	// ถ้า value(ค่าของข้อมูล) ไม่ว่างให้แปลงจาก json เป็น struct เพื่อนำไปใช้ต่อ ปร๊ะ ๆ  ๆ ๆ ๆ ๆ//
	if value != nil {
		if err := json.Unmarshal(b, value); err != nil {
			return err
		}
	}

	return nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.MongoClient.Del(ctx, key).Err()
}

func (r *Redis) DeletePattern(ctx context.Context, key string) error {
	keys, err := r.MongoClient.Keys(ctx, key).Result()
	if err != nil {
		return err
	}
	return r.MongoClient.Del(ctx, keys...).Err()
}

func (r *Redis) IsKeyNotFound(err error) bool {
	return err == redis.Nil
}
