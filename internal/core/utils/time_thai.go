package utils

import (
	"fmt"
	"log"
	"time"
)

func ConvertTimeResponse(createdAt time.Time) string {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return ""
	}
	t := createdAt.In(loc)

	thaiMonths := [...]string{
		"", "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.",
		"ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค.",
	}

	day := t.Day()
	month := thaiMonths[int(t.Month())]
	year := t.Year()

	return fmt.Sprintf("%d %s %d", day, month, year)
}
