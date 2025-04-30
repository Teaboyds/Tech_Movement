package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadedFile struct {
	OriginalName string // ชื่อไฟล์เดิม
	SavedName    string // ชื่อ UUID ที่บันทึกลง disk
}

func UploadFile(c *fiber.Ctx, fieldName string, maxSize int64, uploadDir string) (*UploadedFile, error) {
	fileHeader, err := c.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	if fileHeader.Size > maxSize {
		return nil, fmt.Errorf("file too large. Max: %d bytes", maxSize)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return nil, fmt.Errorf("invalid file type")
	}

	originalName := fileHeader.Filename

	newFileName := uuid.New().String() + ext
	savePath := filepath.Join(uploadDir, newFileName)

	fmt.Printf("savePath: %v\n", savePath)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("cannot create folder: %v", err)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}
	defer src.Close()

	out, err := os.Create(savePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create file: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return nil, fmt.Errorf("cannot save file: %v", err)
	}

	return &UploadedFile{
		OriginalName: originalName,
		SavedName:    newFileName,
	}, nil
}

func AttachBaseURLToImage(filetype string, path string) string {
	baseURL := os.Getenv("BASE_URL")
	RealURL := baseURL + filetype + "/" + path
	return RealURL
}

func DeleteFileInLocalStorage(fileType string, name string) error {
	fullPath := "../upload/" + fileType + "/" + name
	fmt.Printf("fullPath: %q\n", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("Mongo : Image not found at: %s", fullPath)
		return nil
	}

	info, err := os.Stat(fullPath)
	if err == nil {
		fmt.Printf("File Mode: %v\n", info.Mode())
	}

	err = os.Remove(fullPath)
	if err != nil {
		log.Printf("Failed to delete image at %s: %v", fullPath, err)
		return err
	}

	log.Printf("Image deleted: %s", fullPath)
	return nil
}
