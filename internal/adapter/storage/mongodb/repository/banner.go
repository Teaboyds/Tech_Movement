package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	mongoUtils "backend_tech_movement_hex/internal/adapter/storage/mongodb/utils"
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoBannerRepository struct {
	collection *mongo.Collection
}

func NewBannersRepoMongo(db *mongodb.Database) port.BannerRepository {
	return &MongoBannerRepository{
		collection: db.Collection("banner"),
	}
}

func (n *MongoBannerRepository) SaveBanner( /*parameter*/ banner *domain.Banner) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	CateOBJ, err := mongoUtils.ConvertStringToObjectID(banner.Category)
	if err != nil {
		return err
	}

	var objIDs []primitive.ObjectID
	for _, id := range banner.Img {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objIDs = append(objIDs, objID)
	}

	bannerDoc := &models.MongoBanner{
		Title:       banner.Title,
		ContentType: banner.ContentType,
		Status:      banner.Status,
		CategoryID:  CateOBJ,
		ImageID:     objIDs,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = n.collection.InsertOne(ctx, bannerDoc)
	return err
}

func (n *MongoBannerRepository) Retrive(id string) (*domain.Banner, error) {

	ObjID, err := mongoUtils.ConvertStringToObjectID(id)
	if err != nil {
		return nil, err
	}

	var banner models.MongoBanner

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	filter := bson.M{"_id": ObjID}

	err = n.collection.FindOne(ctx, filter).Decode(&banner)
	if err != nil {
		return nil, err
	}

	response := &domain.Banner{
		ID:          banner.ID.Hex(),
		Title:       banner.Title,
		ContentType: banner.ContentType,
		Status:      banner.Status,
		Category:    banner.CategoryID.Hex(),
	}

	return response, nil
}
