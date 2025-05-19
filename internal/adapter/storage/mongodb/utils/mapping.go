package mongoUtils

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	"backend_tech_movement_hex/internal/core/domain"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapNewsToMongo(news *domain.News) (*models.MongoNews, error) {
	// แปลง categoryID จาก string เป็น primitive.ObjectID
	categoryObjectID, err := primitive.ObjectIDFromHex(news.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	thumnailID, err := primitive.ObjectIDFromHex(news.ThumnailID)
	if err != nil {
		return nil, fmt.Errorf("invalid Thumanil ID: %w", err)
	}

	// แปลง image IDs จาก string เป็น primitive.ObjectID
	var imageObjectIDs []primitive.ObjectID
	for _, id := range news.ImageIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("error converting image ID: %w", err)
		}
		imageObjectIDs = append(imageObjectIDs, objID)
	}

	bools, err := strconv.ParseBool(news.Status)
	if err != nil {
		return nil, err
	}

	parseView, err := strconv.Atoi(news.View)
	if err != nil {
		return nil, err
	}

	// สร้าง MongoNews จาก domain.News
	mongoNews := &models.MongoNews{
		ThumnailID:  thumnailID,
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		ImageIDs:    imageObjectIDs,
		CategoryID:  categoryObjectID,
		Tags:        news.Tags,
		Status:      bools,
		ContentType: news.ContentType,
		View:        parseView,
		CreatedAt:   time.Now(), // หรือใช้เวลาจาก news ถ้ามี
		UpdatedAt:   time.Now(), // หรือใช้เวลาจาก news ถ้ามี
	}

	// คืนค่าผลลัพธ์
	return mongoNews, nil
}
