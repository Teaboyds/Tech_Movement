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
	RedisClient *redis.Client
}

func ConnectedRedis(ctx context.Context, config *config.Redis) (port.CacheRepository, error) {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_PORT,
		Password: config.REDIS_PASSWORD,
		DB:       config.REDIS_DB,
	})
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Cannot connect to Redis: ", err)
	}
	log.Println("Connected to Redis")
	return &Redis{RedisClient}, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, timeout time.Duration) error {

	// แปลง value เพื่อเก็บใน redis  //
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err = r.RedisClient.Set(ctx, key, j, timeout).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {

	// ดึง keys มาแล้วแปลงเป็น byte //
	b, err := r.RedisClient.Get(ctx, key).Bytes()
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
	return r.RedisClient.Del(ctx, key).Err()
}

func (r *Redis) DeletePattern(ctx context.Context, key string) error {
	keys, err := r.RedisClient.Keys(ctx, key).Result()
	if err != nil {
		return err
	}
	return r.RedisClient.Del(ctx, keys...).Err()
}

func (r *Redis) IncrementVersion(ctx context.Context, key string) (int64, error) {
	val, err := r.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (r *Redis) IsKeyNotFound(err error) bool {
	return err == redis.Nil
}
