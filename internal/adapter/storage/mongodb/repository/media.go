package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	"backend_tech_movement_hex/internal/core/domain"
	d "backend_tech_movement_hex/internal/core/domain"
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

type MongoMediaRepository struct {
	db *mongo.Collection
}

func NewMediaRepositoryMongo(db *mongodb.Database) port.MediaRepository {
	return &MongoMediaRepository{
		db: db.Collection("media"),
	}
}

func (med *MongoMediaRepository) CreateMedia(media *domain.Media) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	cateObj, err := primitive.ObjectIDFromHex(media.CategoryID)
	if err != nil {
		return fmt.Errorf("cate_obj error in media repo")
	}

	mediaDoc := &models.MongoMedia{
		Title:      media.Title,
		Content:    media.Content,
		URL:        media.URL,
		CategoryID: cateObj,
		Status:     media.Status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = med.db.InsertOne(ctx, mediaDoc)
	if err != nil {
		return fmt.Errorf("cannot save media in mongodb : %s", err)
	}

	return err
}

func (med *MongoMediaRepository) GetVideoHome(cateId string) ([]*domain.Media, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	categoryObjectID, err := primitive.ObjectIDFromHex(cateId)
	if err != nil {
		log.Println("invalid ObjectID from category ID:", err)
		return nil, err
	}

	fmt.Printf("categoryObjectID: %v\n", categoryObjectID)

	var results []models.MongoMedia

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "status", Value: true},
			{Key: "category_id", Value: bson.M{"$ne": categoryObjectID}},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "created_at", Value: -1},
		}}},
		{{Key: "$limit", Value: 4}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "categories"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "category"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$category"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := med.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var responses []*domain.Media
	for _, result := range results {

		responses = append(responses, &d.Media{
			Title:      result.Title,
			Content:    result.Content,
			URL:        result.URL,
			CategoryID: result.CategoryID.Hex(),
			CreatedAt:  utils.ConvertTimeResponse(result.CreatedAt),
		})
	}

	fmt.Printf("responses: %v\n", responses)

	return responses, nil
}

func (med *MongoMediaRepository) GetShortVideoHome(cateId string) ([]*domain.Media, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var lastNews []models.MongoMedia

	categoryObjectID, err := primitive.ObjectIDFromHex(cateId)
	if err != nil {
		log.Println("invalid ObjectID from category ID:", err)
		return nil, err
	}

	findOptions := options.Find().
		SetProjection(bson.M{
			"_id":         0,
			"content":     0,
			"category_id": 0,
			"created_at":  0,
			"updated_at":  0,
			"status":      0,
		}).
		SetLimit(4).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	filter := bson.M{
		"status":      true,
		"category_id": categoryObjectID,
	}

	cursor, err := med.db.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &lastNews); err != nil {
		log.Printf("Error decoding repo last news: %v", err)
		return nil, err
	}
	var shortVideos []*domain.Media
	for _, media := range lastNews {
		shortVideos = append(shortVideos, &domain.Media{
			Title: media.Title,
			URL:   media.URL,
		})
	}

	fmt.Printf("shortVideos: %v\n", shortVideos)

	return shortVideos, nil
}
