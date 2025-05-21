package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	mongoMapper "backend_tech_movement_hex/internal/adapter/storage/mongodb/mapper"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	mongoUtils "backend_tech_movement_hex/internal/adapter/storage/mongodb/utils"
	dif "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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

	saveInfo, err := mongoMapper.InfographicDomainToMongo(info)
	if err != nil {
		return err
	}

	_, err = ip.collection.InsertOne(ctx, saveInfo)
	return err
}

func (ip *MongoInfographicRepository) Retrive(id string) (*dif.Infographic, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjID, err := mongoUtils.ConvertStringToObjectID(id)
	if err != nil {
		return nil, err
	}

	var infographic models.MongoInfographic
	filter := bson.M{"_id": ObjID}

	err = ip.collection.FindOne(ctx, filter).Decode(&infographic)
	if err != nil {
		return nil, bson.ErrDecodeToNil
	}

	infographicResp := mongoMapper.InfographicMongoToDomain(infographic)

	return infographicResp, err
}

func (ip *MongoInfographicRepository) RetrivesInfographic(cateId, sort, view, limit, page string) ([]*dif.Infographic, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	finder := bson.M{}

	parseLimit, err := parseInt(limit)
	if err != nil {
		return nil, err
	}
	parsePage, err := parseInt(page)
	if err != nil {
		return nil, err
	}

	if cateId != "" {
		objID, err := mongoUtils.ConvertStringToObjectID(cateId)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		finder["category"] = objID
	}

	opts := options.Find()

	if sort == "newest" {
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	} else if sort == "oldest" {
		opts.SetSort(bson.D{{Key: "created_at", Value: 1}})
	}

	if view == "asc" {
		opts.SetSort(bson.D{{Key: "view", Value: -1}})
	} else if view == "desc" {
		opts.SetSort(bson.D{{Key: "view", Value: 1}})
	}

	fmt.Printf("view: %v\n", view)

	skip := (parsePage - 1) * parseLimit
	opts.SetLimit(parseLimit).SetSkip(skip)

	cursor, err := ip.collection.Find(ctx, finder, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var InfoResp []*dif.Infographic
	for cursor.Next(ctx) {
		var infographic models.MongoInfographic
		if err := cursor.Decode(&infographic); err != nil {
			return nil, err
		}

		mediaDTO := mongoMapper.InfographicMongoToDomain(infographic)

		InfoResp = append(InfoResp, mediaDTO)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("InfoResp: %v\n", InfoResp)

	return InfoResp, err
}

// func (ip *MongoInfographicRepository) GetInfoHome() ([]*dif.Infographic, error) {
// 	ctx, cancel := utils.NewTimeoutContext()
// 	defer cancel()

// 	var results []models.MongoInfographic

// 	filter := bson.M{}
// 	findOptions := options.Find().
// 		SetLimit(5).
// 		SetSort(bson.D{{Key: "_id", Value: -1}})

// 	cursor, err := ip.collection.Find(ctx, filter, findOptions)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot aggregate: %s", err)
// 	}
// 	defer cursor.Close(ctx)

// 	if err := cursor.All(ctx, &results); err != nil {
// 		return nil, fmt.Errorf("cannot cursor info: %s", err)
// 	}

// 	var Infographic []*domain.Infographic
// 	for _, short := range results {

// 		Infographic = append(Infographic, &domain.Infographic{
// 			ID:        short.ID.Hex(),
// 			Title:     short.Title,
// 			Category:  short.Category.Hex(),
// 			Tags:      short.Tags,
// 			Image:     short.Image.Hex(),
// 			Status:    short.Status,
// 			CreatedAt: short.CreatedAt.String(),
// 			UpdatedAt: short.UpdatedAt.String(),
// 		})
// 	}

// 	fmt.Printf("shortVideos: %v\n", Infographic)
// 	return Infographic, nil
// }
