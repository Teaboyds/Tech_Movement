package utils

import (
	"backend_tech_movement_hex/internal/core/domain"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

func ProcessImageToURL(news *domain.News) error {
	// Decode Base64 string
	data, err := base64.StdEncoding.DecodeString(news.Image)
	if err != nil {
		return fmt.Errorf("failed to decode Base64: %v", err)
	}

	// กำหนดพาธไฟล์
	uploadPath := "./internal/core/upload"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, os.ModePerm)
	}

	// สร้างชื่อไฟล์ (อาจใช่ ID หรือชื่อเฉพาะ)
	fileName := fmt.Sprintf("%s.png", news.ID.Hex())
	filePath := filepath.Join(uploadPath, fileName)

	// เขียนข้อมูลลงไฟล์
	if err := os.WriteFile(filePath, data, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	// สร้าง URL สำหรับไฟล์
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5050/api/v1/news/uploads"
	}
	news.Image = fmt.Sprintf("%s/%s", baseURL, fileName)

	return nil
}
