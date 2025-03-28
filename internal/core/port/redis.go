package port

import (
	"context"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, timeout time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, key string) error
	IsKeyNotFound(err error) bool
}
