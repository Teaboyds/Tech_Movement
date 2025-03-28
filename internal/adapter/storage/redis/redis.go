package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, timeout time.Duration) error {

	// แปลง value เพื่อเก็บใน redis เพราะ ฝ //
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err = r.client.Set(ctx, key, j, timeout).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {

	// ดึง keys มาแล้วแปลงเป็น byte //
	b, err := r.client.Get(ctx, key).Bytes()
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
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) DeletePattern(ctx context.Context, key string) error {
	keys, err := r.client.Keys(ctx, key).Result()
	if err != nil {
		return err
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *Redis) IsKeyNotFound(err error) bool {
	return err == redis.Nil
}
