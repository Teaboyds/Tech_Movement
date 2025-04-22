package utils

import (
	"backend_tech_movement_hex/internal/core/domain"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

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

func UploadFile(c *fiber.Ctx, fieldName string, maxSize int64, uploadDir string) (string, error) {
	// รับไฟล์จากฟอร์ม
	fileHeader, err := c.FormFile(fieldName)
	if err != nil {
		log.Println("Error receiving file:", err)
		return "", err
	}

	if fileHeader.Size > maxSize {
		return "", fmt.Errorf("file too large. Maximum size is %d bytes", maxSize)
	}

	fileEXT := filepath.Ext(fileHeader.Filename)
	if fileEXT != ".png" && fileEXT != ".jpeg" && fileEXT != ".jpg" {
		return "", fmt.Errorf("invalid file type. Only .png, .jpeg, and .jpg are allowed")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open file: %v", err)
	}
	defer src.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", fmt.Errorf("failed to generate file hash: %v", err)
	}
	hashString := hex.EncodeToString(hash.Sum(nil))
	fileName := hashString + fileEXT

	savePath := filepath.Join(uploadDir, fileName)
	if _, err := os.Stat(savePath); os.IsNotExist(err) {

		srcAgain, err := fileHeader.Open()
		if err != nil {
			return "", fmt.Errorf("cannot reopen file: %v", err)
		}
		defer srcAgain.Close()

		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("cannot create folder: %v", err)
		}

		out, err := os.Create(savePath)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %v", err)
		}
		defer out.Close()

		if _, err := io.Copy(out, srcAgain); err != nil {
			return "", fmt.Errorf("failed to save file: %v", err)
		}
	}

	return fileName, nil
}

func AttachBaseURLToImage(new *domain.News) {
	baseURL := os.Getenv("BASE_URL")
	new.Image = baseURL + new.Image
}
