package mongoMapper

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	"backend_tech_movement_hex/internal/core/domain"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InfographicDomainToMongo(info *domain.Infographic) (*models.MongoInfographic, error) {

	CateOBJ, err := primitive.ObjectIDFromHex(info.Category)
	if err != nil {
		log.Println("err 1 ", err)
		return nil, err
	}

	FileOBJ, err := primitive.ObjectIDFromHex(info.Image)
	if err != nil {
		log.Println("err 2 ", err)
		return nil, err
	}

	parseStatus, err := strconv.ParseBool(info.Status)
	if err != nil {
		return nil, err
	}

	parsePageView, err := strconv.Atoi(info.PageView)
	if err != nil {
		return nil, err
	}

	return &models.MongoInfographic{
		Title:     info.Title,
		Image:     FileOBJ,
		Category:  CateOBJ,
		Tags:      info.Tags,
		Status:    parseStatus,
		PageView:  parsePageView,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, err
}

func InfographicMongoToDomain(info models.MongoInfographic) *domain.Infographic {

	parseStatus := strconv.FormatBool(info.Status)

	return &domain.Infographic{
		ID:        info.ID.Hex(),
		Image:     info.Image.Hex(),
		Title:     info.Title,
		Category:  info.Category.Hex(),
		Tags:      info.Tags,
		Status:    parseStatus,
		CreatedAt: info.CreatedAt.Format(time.RFC3339),
		UpdatedAt: info.UpdatedAt.Format(time.RFC3339),
	}
}
