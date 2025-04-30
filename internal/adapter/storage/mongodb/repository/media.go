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
	db           *mongo.Collection
	categoryRepo port.CategoryRepository
}

func NewMediaRepositoryMongo(db *mongodb.Database, categoryRepo port.CategoryRepository) port.MediaRepository {
	return &MongoMediaRepository{
		db:           db.Collection("media"),
		categoryRepo: categoryRepo,
	}
}

func (med *MongoMediaRepository) CreateMedia(media *domain.MediaRequest) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	Category, err := med.categoryRepo.GetByID(media.Category)
	if err != nil {
		return err
	}

	CateOBJ, err := primitive.ObjectIDFromHex(Category.ID)
	if err != nil {
		return fmt.Errorf("cate_obj error in media repo")
	}

	fmt.Printf("CateOBJ: %v\n", CateOBJ)

	medDoc := &models.MongoMedia{
		Title:      media.Title,
		Content:    media.Content,
		URL:        media.URL,
		CategoryID: CateOBJ,
		Status:     media.Status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = med.db.InsertOne(ctx, medDoc)

	return err
}
func (med *MongoMediaRepository) GetVideoHome() ([]*domain.VideoResponse, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	category, err := med.categoryRepo.GetByName("Short Video")
	if err != nil {
		log.Println("error fetching short_video category:", err)
		return nil, err
	}

	categoryObjectID, err := primitive.ObjectIDFromHex(category.ID)
	if err != nil {
		log.Println("invalid ObjectID from category ID:", err)
		return nil, err
	}

	log.Printf("Excluded category ID: %v (type: %T)", category.ID, category.ID)

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
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category_id"},
			{Key: "news", Value: bson.D{
				{Key: "$first", Value: "$$ROOT"},
			}},
		}}},
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

	// Execute the aggregation pipeline
	cursor, err := med.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Prepare the results struct to hold the aggregation response
	var results []struct {
		Media    models.MongoMedia    `bson:"news"`
		Category models.MongoCategory `bson:"category"`
	}

	// Decode the aggregation results
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Map the results to a slice of VideoResponse
	var responses []*domain.VideoResponse
	for _, result := range results {
		categoryResp := d.CategoryResponse{
			ID:   result.Category.ID.Hex(),
			Name: result.Category.Name,
		}

		responses = append(responses, &d.VideoResponse{
			Title:     result.Media.Title,
			Content:   result.Media.Content,
			URL:       result.Media.URL,
			Category:  categoryResp,
			CreatedAt: utils.ConvertTimeResponse(result.Media.CreatedAt),
		})
	}

	fmt.Printf("responses: %v\n", responses)

	return responses, nil
}

func (med *MongoMediaRepository) GetShortVideoHome() ([]*domain.ShortVideo, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var lastNews []models.MongoMedia

	categories, err := med.categoryRepo.GetByName("Short Video")
	if err != nil {
		return nil, fmt.Errorf("error fetching category: %w", err)
	}

	categoryObjectID, err := primitive.ObjectIDFromHex(categories.ID)
	if err != nil {
		log.Println("invalid ObjectID from category ID:", err)
		return nil, err
	}

	if categories == nil {
		return nil, fmt.Errorf("category 'short vdo' not found")
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
	var shortVideos []*domain.ShortVideo
	for _, media := range lastNews {
		shortVideos = append(shortVideos, &domain.ShortVideo{
			Title: media.Title,
			URL:   media.URL,
		})
	}

	fmt.Printf("shortVideos: %v\n", shortVideos)

	return shortVideos, nil
}
