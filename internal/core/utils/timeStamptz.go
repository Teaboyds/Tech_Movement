package utils

import (
	"time"
)

func SetTimestamps() (string, string) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc).Format(time.RFC3339)
	return now, now
}
