package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	mongoMapper "backend_tech_movement_hex/internal/adapter/storage/mongodb/mapper"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	mongoUtils "backend_tech_movement_hex/internal/adapter/storage/mongodb/utils"
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (n *MongoBannerRepository) SaveBanner(banner *domain.Banner) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	bannerDoc := &models.MongoBanner{
		DesktopImage: models.ImageInfo(banner.DesktopImage),
		MobileImage:  models.ImageInfo(banner.MobileImage),
		Status:       models.StatusType(banner.Status),
		LinkUrl:      banner.LinkUrl,
		Action:       banner.Action,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := n.collection.InsertOne(ctx, bannerDoc)
	return err
}

func (n *MongoBannerRepository) Retrive(id string) (*domain.Banner, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	objID, err := mongoUtils.ConvertStringToObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("cannot convert string to objected id : %v", err)
	}

	var banner models.MongoBanner

	filter := bson.M{"_id": objID}

	err = n.collection.FindOne(ctx, filter).Decode(&banner)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch banner info : %v", err)
	}

	bannerResponse := mongoMapper.BannerToDomain(banner)

	return bannerResponse, err
}

func (n *MongoBannerRepository) Retrives(page_type string) ([]*domain.Banner, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	finder := bson.M{}

	if page_type != "" {
		finder["status."+page_type] = true

	}

	cursor, err := n.collection.Find(ctx, finder)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bannersResp []*domain.Banner
	for cursor.Next(ctx) {
		var banners models.MongoBanner
		if err := cursor.Decode(&banners); err != nil {
			return nil, err
		}

		bannersDTO := mongoMapper.BannerToDomain(banners)

		bannersResp = append(bannersResp, bannersDTO)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bannersResp, err
}

func (n *MongoBannerRepository) SaveBannerV2(banner *domain.BannerV2) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	bannerDoc := &models.MongoBannerV2{
		DesktopImage: banner.DesktopImage,
		MobileImage:  banner.MobileImage,
		Status:       models.StatusType(banner.Status),
		LinkUrl:      banner.LinkUrl,
		Action:       banner.Action,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := n.collection.InsertOne(ctx, bannerDoc)
	return err
}

func (n *MongoBannerRepository) RetriveV2(id string) (*domain.BannerV2, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	objID, err := mongoUtils.ConvertStringToObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("cannot convert string to objected id : %v", err)
	}

	var banner models.MongoBannerV2

	filter := bson.M{"_id": objID}

	err = n.collection.FindOne(ctx, filter).Decode(&banner)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch banner info : %v", err)
	}

	bannerResponse := mongoMapper.MongoBannerToDomain(banner)

	return bannerResponse, err
}
