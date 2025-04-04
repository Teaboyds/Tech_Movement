package utils

import (
	"image"
	"os"

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
