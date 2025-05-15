package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	mongoUtils "backend_tech_movement_hex/internal/adapter/storage/mongodb/utils"
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (n *MongoBannerRepository) Retrives() ([]*domain.Banner, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var banners []*models.MongoBanner

	findBanner := options.Find().
		SetLimit(5).
		SetSort(bson.D{{"created_at", -1}})

	filter := bson.M{
		"status": true,
	}
	cursor, err := n.collection.Find(ctx, filter, findBanner)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &banners); err != nil {
		return nil, err
	}

	var responses []*domain.Banner
	for _, banner := range banners {
		var imageIDs []string
		for _, oid := range banner.ImageID {
			imageIDs = append(imageIDs, oid.Hex())
		}

		responseBanner := &domain.Banner{
			ID:          banner.ID.Hex(),
			Title:       banner.Title,
			ContentType: banner.ContentType,
			Status:      banner.Status,
			Category:    banner.CategoryID.Hex(),
			Img:         imageIDs,
		}
		responses = append(responses, responseBanner)
	}
	return responses, nil
}

func (n *MongoBannerRepository) Updated(id string, banner *domain.Banner) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	CateObj, err := primitive.ObjectIDFromHex(banner.Category)
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
	fmt.Printf("objIDs: %v\n", objIDs)

	update := bson.M{
		"$set": bson.M{
			"title":        banner.Title,
			"content_type": banner.ContentType,
			"status":       banner.Status,
			"category_id":  CateObj,
			"image_id":     objIDs,
			"updated_at":   time.Now(),
		},
	}

	filter := bson.M{"_id": ObjID}
	result, err := n.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (banner *MongoBannerRepository) Delete(id string) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := banner.collection.DeleteOne(ctx, bson.M{"_id": ObjID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
