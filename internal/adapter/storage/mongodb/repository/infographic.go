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
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoInfographicRepository struct {
	collection   *mongo.Collection
	categoryRepo port.CategoryRepository
	fileRepo     port.UploadRepository
}

func NewInfographicRepositoryMongo(
	db *mongodb.Database,
	categoryRepo port.CategoryRepository,
	fileRepo port.UploadRepository) port.InfographicRepository {
	return &MongoInfographicRepository{
		collection:   db.Collection("infographic"),
		categoryRepo: categoryRepo,
		fileRepo:     fileRepo,
	}
}

func (ip *MongoInfographicRepository) CreateInfo(info *dif.InfographicRequest) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	Category, err := ip.categoryRepo.GetByID(info.Category)
	if err != nil {
		return err
	}

	fmt.Printf("Category: %v\n", Category)

	File, err := ip.fileRepo.GetFileByID(info.Image)
	if err != nil {
		return err
	}

	fmt.Printf("File: %v\n", File)

	CateOBJ, err := primitive.ObjectIDFromHex(Category.ID)
	if err != nil {
		log.Println("err 1 ", err)
		return err
	}

	FileOBJ, err := primitive.ObjectIDFromHex(File.ID)
	if err != nil {
		log.Println("err 2 ", err)
		return err
	}

	statusBool, err := strconv.ParseBool(info.Status)
	if err != nil {
		return err
	}

	newInfo := &models.MongoInfographic{
		Title:     info.Title,
		Image:     FileOBJ,
		Category:  CateOBJ,
		Tags:      info.Tags,
		Status:    statusBool,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = ip.collection.InsertOne(ctx, newInfo)
	return err
}

func (ip *MongoInfographicRepository) GetInfoHome() ([]dif.InfographicRespose, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"status": true}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: -1}}}},
		{{Key: "$limit", Value: 6}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "file_Upload",
			"localField":   "image",
			"foreignField": "_id",
			"as":           "image_info",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$image_info",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":        1,
			"title":      1,
			"created_at": 1,
			"image_info": 1,
		}}},
	}

	cursor, err := ip.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID        primitive.ObjectID           `bson:"_id"`
		Title     string                       `bson:"title"`
		CreatedAt time.Time                    `bson:"created_at"`
		ImageInfo models.MongoUploadRepository `bson:"image_info"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Map to response
	response := make([]domain.InfographicRespose, len(results))
	for i, r := range results {
		response[i] = domain.InfographicRespose{
			ID:    r.ID.Hex(),
			Title: r.Title,
			Image: domain.UploadFileResponse{
				ID:       r.ImageInfo.ID.Hex(),
				Path:     r.ImageInfo.Path,
				Name:     r.ImageInfo.Name,
				FileType: r.ImageInfo.FileType,
			},
			CreatedAt: utils.ConvertTimeResponse(r.CreatedAt),
		}
	}

	return response, nil
}
