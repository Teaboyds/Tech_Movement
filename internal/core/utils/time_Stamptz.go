package utils

import (
	"time"
)

func SetTimestamps() (string, string) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc).Format("2006-01-02 15:04:05.00")
	return now, now
}
