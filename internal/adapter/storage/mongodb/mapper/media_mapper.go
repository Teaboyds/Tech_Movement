package mongoMapper

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	mongoUtils "backend_tech_movement_hex/internal/adapter/storage/mongodb/utils"
	"backend_tech_movement_hex/internal/core/domain"
	"strconv"
	"time"
)

func MediaDomainToMongo(media domain.Media) models.MongoMedia {

	parseObj, _ := mongoUtils.ConvertStringToObjectID(media.Thumnail)

	parseCateObj, _ := mongoUtils.ConvertStringToObjectID(media.Category)

	parseView, _ := strconv.Atoi(media.View)

	return models.MongoMedia{
		Title:     media.Title,
		Content:   media.Content,
		VideoUrl:  media.VideoURL,
		Thumnail:  parseObj,
		Category:  parseCateObj,
		Tags:      media.Tags,
		View:      parseView,
		Action:    media.Action,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func MediaMongoToDomain(media models.MongoMedia) *domain.Media {

	parseView := strconv.Itoa(media.View)

	return &domain.Media{
		ID:        media.ID.Hex(),
		Title:     media.Title,
		Content:   media.Content,
		VideoURL:  media.VideoUrl,
		Thumnail:  media.Thumnail.Hex(),
		Category:  media.Category.Hex(),
		Tags:      media.Tags,
		View:      parseView,
		Action:    media.Action,
		CreatedAt: media.CreatedAt.Format(time.RFC3339),
		UpdatedAt: media.UpdatedAt.Format(time.RFC3339),
	}
}
