package utils

import (
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func GetImageDimensions(filePath string) (int, int) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0
	}
	return img.Width, img.Height
}

func GenerateUUID() string {
	return uuid.New().String()
}

func DeleteImageFiles(fileName string) {
	if fileName == "" {
		return
	}

	path := filepath.Join("./upload/image", fileName)

	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to delete file: %s, error: %v", path, err)
	}
}
