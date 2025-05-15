package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	"backend_tech_movement_hex/internal/core/domain"
	dif "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInfographicRepository struct {
	collection *mongo.Collection
}

func NewInfographicRepositoryMongo(
	db *mongodb.Database) port.InfographicRepository {
	return &MongoInfographicRepository{
		collection: db.Collection("infographic"),
	}
}

func (ip *MongoInfographicRepository) CreateInfo(info *dif.Infographic) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	fmt.Printf("Category: %v\n", info.Category)

	CateOBJ, err := primitive.ObjectIDFromHex(info.Category)
	if err != nil {
		log.Println("err 1 ", err)
		return err
	}

	FileOBJ, err := primitive.ObjectIDFromHex(info.Image)
	if err != nil {
		log.Println("err 2 ", err)
		return err
	}

	newInfo := &models.MongoInfographic{
		Title:     info.Title,
		Image:     FileOBJ,
		Category:  CateOBJ,
		Tags:      info.Tags,
		Status:    info.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = ip.collection.InsertOne(ctx, newInfo)
	return err
}

func (ip *MongoInfographicRepository) GetInfoHome() ([]*dif.Infographic, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var results []models.MongoInfographic

	filter := bson.M{}
	findOptions := options.Find().
		SetLimit(5).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := ip.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("cannot aggregate: %s", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("cannot cursor info: %s", err)
	}

	var Infographic []*domain.Infographic
	for _, short := range results {

		Infographic = append(Infographic, &domain.Infographic{
			ID:        short.ID.Hex(),
			Title:     short.Title,
			Category:  short.Category.Hex(),
			Tags:      short.Tags,
			Image:     short.Image.Hex(),
			Status:    short.Status,
			CreatedAt: short.CreatedAt.String(),
			UpdatedAt: short.UpdatedAt.String(),
		})
	}

	fmt.Printf("shortVideos: %v\n", Infographic)
	return Infographic, nil
}
