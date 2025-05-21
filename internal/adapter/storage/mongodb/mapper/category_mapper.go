package mongoMapper

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	mongoUtils "backend_tech_movement_hex/internal/adapter/storage/mongodb/utils"
	"backend_tech_movement_hex/internal/core/domain"
	"time"
)

func CategoryDomainToMongo(category domain.Category) *models.MongoCategory {

	parseObj, err := mongoUtils.ConvertStringToObjectID(category.ID)
	if err != nil {
		return nil
	}

	return &models.MongoCategory{
		ID:           parseObj,
		Name:         category.Name,
		CategoryType: category.CategoryType,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func CategoryMongoToDomain(category models.MongoCategory) *domain.Category {
	return &domain.Category{
		ID:           category.ID.Hex(),
		Name:         category.Name,
		CategoryType: category.CategoryType,
		CreatedAt:    category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    category.UpdatedAt.Format(time.RFC3339),
	}
}
