// secondary adapters //
package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	mongoUtils "backend_tech_movement_hex/internal/adapter/storage/mongodb/utils"
	"backend_tech_movement_hex/internal/core/domain"
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoNewsRepository struct {
	collection *mongo.Collection
}

func NewNewsRepo(db *mongodb.Database) port.NewsRepository {
	return &MongoNewsRepository{
		collection: db.Collection("news"),
	}
}

// // index model ////
func (n *MongoNewsRepository) EnsureNewsIndexs() error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "_id", Value: -1},
		},
		Options: options.Index().SetName("status_id_desc"),
	}

	_, err := n.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

// /// create area ///
func (n *MongoNewsRepository) SaveNews(news *d.News) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	mongoNews, err := mongoUtils.MapNewsToMongo(news)
	if err != nil {
		return err
	}

	_, err = n.collection.InsertOne(ctx, mongoNews)
	return err
}

// /// get area ///

func (n *MongoNewsRepository) GetNewsByID(id string) (*d.News, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("objID: %v\n", objID)

	var mongoNews models.MongoNews
	err = n.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&mongoNews)
	if err != nil {
		return nil, err
	}

	var imageIDs []string
	for _, oid := range mongoNews.ImageIDs {
		imageIDs = append(imageIDs, oid.Hex())
	}

	parseView := strconv.Itoa(mongoNews.View)

	response := &domain.News{
		ID:          mongoNews.ID.Hex(),
		ThumnailID:  mongoNews.ThumnailID.Hex(),
		Title:       mongoNews.Title,
		Description: mongoNews.Description,
		Content:     mongoNews.Content,
		ImageIDs:    imageIDs,
		CategoryID:  mongoNews.CategoryID.Hex(),
		Tags:        mongoNews.Tags,
		Status:      strconv.FormatBool(mongoNews.Status),
		ContentType: mongoNews.ContentType,
		View:        parseView,
		CreatedAt:   mongoNews.CreatedAt.Format(time.RFC1123),
		UpdatedAt:   mongoNews.UpdatedAt.Format(time.RFC1123),
	}

	return response, nil
}

func (n *MongoNewsRepository) GetLastNews() ([]*d.News, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var lastNews []models.MongoNews

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "status", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "created_at", Value: 1},
		}}},
		{{Key: "$limit", Value: 5}},
	}

	cursor, err := n.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &lastNews); err != nil {
		return nil, err
	}

	var responseNews []*domain.News
	for _, news := range lastNews {
		imageIDs := make([]string, 0, len(news.ThumnailID))
		for _, oid := range news.ImageIDs {
			imageIDs = append(imageIDs, oid.Hex())
		}

		parseView := strconv.Itoa(news.View)

		resp := &domain.News{
			ID:          news.ID.Hex(),
			ThumnailID:  news.ThumnailID.Hex(),
			Title:       news.Title,
			Description: news.Description,
			Content:     news.Content,
			ImageIDs:    imageIDs,
			CategoryID:  news.CategoryID.Hex(),
			Tags:        news.Tags,
			Status:      strconv.FormatBool(news.Status),
			View:        parseView,
			ContentType: news.ContentType,
			CreatedAt:   news.CreatedAt.Format(time.RFC3339),
		}
		responseNews = append(responseNews, resp)
	}

	fmt.Printf("responseNews Repo: %v\n", responseNews)

	return responseNews, nil
}

func (n *MongoNewsRepository) GetTechnologyNews() ([]*d.News, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var lastNews []models.MongoNews

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "status", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "categories"},
			{Key: "localField", Value: "category_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "category"},
		}}},
		{{Key: "$unwind", Value: "$category"}},
		{{Key: "$match", Value: bson.D{
			{Key: "category.category_type", Value: "news"},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "category_id", Value: 1},
			{Key: "created_at", Value: -1},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category_id"},
			{Key: "latest_news", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: bson.D{
				{Key: "$mergeObjects", Value: "$latest_news"},
			}},
		}}},
		{{Key: "$limit", Value: 5}},
	}

	cursor, err := n.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &lastNews); err != nil {
		return nil, err
	}

	var responseNews []*domain.News
	for _, news := range lastNews {
		imageIDs := make([]string, 0, len(news.ThumnailID))
		for _, oid := range news.ImageIDs {
			imageIDs = append(imageIDs, oid.Hex())
		}

		parseView := strconv.Itoa(news.View)

		resp := &domain.News{
			ID:          news.ID.Hex(),
			ThumnailID:  news.ThumnailID.Hex(),
			Title:       news.Title,
			Description: news.Description,
			Content:     news.Content,
			ImageIDs:    imageIDs,
			CategoryID:  news.CategoryID.Hex(),
			Tags:        news.Tags,
			Status:      strconv.FormatBool(news.Status),
			View:        parseView,
			ContentType: news.ContentType,
			CreatedAt:   news.CreatedAt.Format(time.RFC3339),
		}
		responseNews = append(responseNews, resp)
	}

	fmt.Printf("responseNews Repo: %v\n", responseNews)

	return responseNews, nil
}

func (n *MongoNewsRepository) Find(catID, ConType, Sort, status, view, search string, limit, page int64) ([]*d.News, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	Finder := bson.M{}

	fmt.Printf("status: %v\n", status)
	if catID != "" {
		objID, err := mongoUtils.ConvertStringToObjectID(catID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		Finder["category_id"] = objID
	}

	if search != "" {
		Finder["$or"] = bson.A{
			bson.M{"title": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if ConType != "" {
		Finder["content_type"] = ConType
	}

	if status != "" {
		tea, err := strconv.ParseBool(status)
		if err != nil {
			return nil, fmt.Errorf("input it's not 'true' or 'false'")
		}
		Finder["status"] = tea
	}

	opts := options.Find()

	if Sort == "newest" {
		opts.SetSort(bson.D{{Key: "created_at", Value: 1}})
	} else if Sort == "oldest" {
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	if view == "asc" {
		opts.SetSort(bson.D{{Key: "view", Value: 1}})
	} else if view == "desc" {
		opts.SetSort(bson.D{{Key: "view", Value: -1}})
	}

	skip := (page - 1) * limit
	opts.SetLimit(limit).SetSkip(skip)

	cursor, err := n.collection.Find(ctx, Finder, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*d.News
	for cursor.Next(ctx) {
		var news models.MongoNews
		if err := cursor.Decode(&news); err != nil {
			return nil, err
		}

		var imageIDs []string
		for _, oid := range news.ImageIDs {
			imageIDs = append(imageIDs, oid.Hex())
		}

		parseViews := strconv.Itoa(news.View)

		results = append(results, &d.News{
			ID:          news.ID.Hex(),
			ThumnailID:  news.ThumnailID.Hex(),
			Title:       news.Title,
			Description: news.Description,
			Content:     news.Content,
			ImageIDs:    imageIDs,
			CategoryID:  news.CategoryID.Hex(),
			Tags:        news.Tags,
			Status:      strconv.FormatBool(news.Status),
			ContentType: news.ContentType,
			View:        parseViews,
			CreatedAt:   news.CreatedAt.String(),
			UpdatedAt:   news.UpdatedAt.String(),
		})
	}

	return results, nil
}

// /// get area ///

func (n *MongoNewsRepository) UpdateNews(id string, news *d.News) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	cateOBJ, err := primitive.ObjectIDFromHex(news.CategoryID)
	if err != nil {
		return err
	}

	var img []primitive.ObjectID
	for _, id := range news.ImageIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		img = append(img, objID)
	}

	update := bson.M{
		"$set": bson.M{
			"thumnail_id":  news.ThumnailID,
			"title":        news.Title,
			"description":  news.Description,
			"content":      news.Content,
			"image_ids":    img,
			"category_id":  cateOBJ,
			"tags":         news.Tags,
			"status":       news.Status,
			"content_type": news.ContentType,
			"updated_at":   time.Now(),
		},
	}

	_, err = n.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		update,
		options.Update(),
	)

	return err
}

func (n *MongoNewsRepository) Delete(id string) error {
	ctx, cancle := utils.NewTimeoutContext()
	defer cancle()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := n.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		log.Printf("No news found with ID: %s", id)
		return errors.New("news not found or already deleted")
	}

	return nil
}

func (n *MongoNewsRepository) DeleteMany(id []string) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var objIDs []primitive.ObjectID
	for _, id := range id {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid image ID format: %w", err)
		}
		objIDs = append(objIDs, objID)
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}

	results, err := n.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("cannot delete news")
	}

	if results.DeletedCount == 0 {
		log.Printf("No news found with ID: %s", id)
		return errors.New("news not found or already deleted")
	}

	return nil
}

func (n *MongoNewsRepository) Count(catID, ConType, Status string) (int64, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	fmt.Printf("Status: %v\n", Status)
	finder := bson.M{}

	if catID != "" {
		objID, err := mongoUtils.ConvertStringToObjectID(catID)
		if err != nil {
			return 0, fmt.Errorf("invalid category ID: %w", err)
		}
		finder["category_id"] = objID
	}

	if ConType != "" {
		finder["content_type"] = ConType
	}

	if Status != "" {
		tea, err := strconv.ParseBool(Status)
		if err != nil {
			return 0, fmt.Errorf("input it's not 'true' or 'false'")
		}
		finder["status"] = tea
	}

	count, err := n.collection.CountDocuments(ctx, finder)
	if err != nil {
		return 0, err
	}

	return count, err
}

// func (n *MongoCategoryRepository) GetNewsByTags(name string) ([]d.News, error) {

// 	var news []d.News
// 	ctx, cancel := utils.NewTimeoutContext()
// 	defer cancel()

// 	filter := bson.M{"tags": name}
// 	cursor, err := n.db.Find(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	if err := cursor.All(ctx, &news); err != nil {
// 		return nil, err
// 	}

// 	return news, nil
// }
