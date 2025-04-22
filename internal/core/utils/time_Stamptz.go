package utils

import (
	"time"
)

func SetTimestamps() (string, string) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc).Format("2006-01-02 15:04:05.00")
	return now, now
}

func SetUpdateTimestamps() string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc).Format("2006-01-02 15:04:05.00")
	return now
}

func ConvertTimestamp(timestampStr string) string {

	loc, _ := time.LoadLocation("Asia/Bangkok")
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05.00", timestampStr, loc)
	if err != nil {
		return "Invalid timestamp format"
	}
	return parsedTime.Format(time.RFC3339)
}
