package mongoUtils

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	"backend_tech_movement_hex/internal/core/domain"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapNewsToMongo(news *domain.News) (*models.MongoNews, error) {
	// แปลง categoryID จาก string เป็น primitive.ObjectID
	categoryObjectID, err := primitive.ObjectIDFromHex(news.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// แปลง image IDs จาก string เป็น primitive.ObjectID
	var imageObjectIDs []primitive.ObjectID
	for _, id := range news.Image {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("error converting image ID: %w", err)
		}
		imageObjectIDs = append(imageObjectIDs, objID)
	}

	// สร้าง MongoNews จาก domain.News
	mongoNews := &models.MongoNews{
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		Image:       imageObjectIDs,
		CategoryID:  categoryObjectID,
		Tag:         news.Tag,
		Status:      news.Status,
		ContentType: news.ContentType,
		CreatedAt:   time.Now(), // หรือใช้เวลาจาก news ถ้ามี
		UpdatedAt:   time.Now(), // หรือใช้เวลาจาก news ถ้ามี
	}

	// คืนค่าผลลัพธ์
	return mongoNews, nil
}
