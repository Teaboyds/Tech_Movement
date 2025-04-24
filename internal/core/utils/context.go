package utils

import (
	"context"
	"time"
)

func NewTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}
