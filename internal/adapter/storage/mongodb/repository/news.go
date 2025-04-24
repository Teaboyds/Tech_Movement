// secondary adapters //
package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"errors"
	"fmt"
	"log"
	"os"
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
	return &MongoNewsRepository{collection: db.Collection("news")}
}

// // index model ////
func (n *MongoNewsRepository) EnsureNewsIndexs() error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "created_at", Value: -1}, // ใหม่ไปเก่า //
			{Key: "category_id._id", Value: -1},
			{Key: "status", Value: -1},
			{Key: "content_status", Value: -1},
			{Key: "content_type", Value: -1},
			{Key: "_id", Value: -1},
		},
	}

	_, err := n.collection.Indexes().CreateOne(ctx, indexModel)

	return err
}

// /// create area ///
func (n *MongoNewsRepository) Create(news *d.News) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	news.ID = primitive.NewObjectID()
	_, err := n.collection.InsertOne(ctx, news)

	return err
}

// /// get area ///

func (n *MongoNewsRepository) GetNewsPagination(lastID string, limit int) ([]d.News, error) {

	var news []d.News
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	// cursor based pagination //
	filter := bson.M{}
	if lastID != "" {
		ObjID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
		filter["_id"] = bson.M{"$lt": ObjID}
	}

	// set spect for sort //
	findOption := options.Find()
	findOption.SetSort(bson.M{"_id": -1})
	findOption.SetLimit(int64(limit))

	cursor, err := n.collection.Find(ctx, filter, findOption)
	if err != nil {
		return nil, err
	}

	// loop decode and appened//
	for cursor.Next(ctx) {
		var new d.News
		if err := cursor.Decode(&new); err != nil {
			return nil, err
		}
		news = append(news, new)
	}
	defer cursor.Close(ctx)

	return news, nil
}

func (n *MongoNewsRepository) GetNewsByID(id string) (*d.News, error) {
	var news d.News
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if n.collection == nil {
		return nil, errors.New("mongoDB client is nil")
	}

	err = n.collection.FindOne(ctx, bson.M{"_id": ObjID}).Decode(&news)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &news, nil
}

func (n *MongoNewsRepository) GetNewsByCategory(categoryID string, lastID string) ([]d.News, string, error) {

	var newsList []d.News
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, "", fmt.Errorf("invalid CategoryID: %v", err)
	}

	filter := bson.M{
		"category_id._id": objID,
		"status":          true,
		"content_status":  "published",
		"content_type":    "general",
	}

	if lastID != "" {
		cursorID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, "", fmt.Errorf("invalid lastID: %v", err)
		}
		filter["_id"] = bson.M{"$lt": cursorID}
	}

	findOptions := options.Find().
		SetProjection(bson.M{"tag": 0, "status": 0, "content_type": 0, "updated_at": 0}).
		SetLimit(9).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := n.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &newsList); err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(newsList) > 0 {
		last := newsList[len(newsList)-1]
		nextCursor = last.ID.Hex()
	}

	return newsList, nextCursor, nil
}

func (n *MongoNewsRepository) GetLastNews() ([]d.News, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var lastNews []d.News

	findOptions := options.Find().
		SetProjection(bson.M{"_id": 0, "tag": 0, "status": 0, "content_type": 0, "updated_at": 0}).
		SetLimit(5).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	filter := bson.M{
		"status":         true,
		"content_status": "published",
		"content_type":   "general",
	}

	cursor, err := n.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &lastNews); err != nil {
		log.Printf("Error decoding repo last news: %v", err)
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return lastNews, nil
}

func (n *MongoNewsRepository) GetNewsByCategoryHomePage(categoryID string) ([]d.News, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var newsCategory []d.News

	filter := bson.M{
		"status":         true,
		"content_status": "published",
		"content_type":   "general",
	}

	if categoryID != "" {
		objID, err := primitive.ObjectIDFromHex(categoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid CategoryID: %v", err)
		}
		filter["category_id._id"] = objID
	}

	findOptions := options.Find().
		SetProjection(bson.M{
			"_id":          0,
			"tag":          0,
			"status":       0,
			"content_type": 0,
			"updated_at":   0,
		}).
		SetLimit(9).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := n.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("cannot fetching news: %v", err)
	}
	if cursor != nil {
		defer cursor.Close(ctx)
	}

	if err := cursor.All(ctx, &newsCategory); err != nil {
		return nil, fmt.Errorf("error decoding news: %v", err)
	}

	return newsCategory, nil
}

func (n *MongoNewsRepository) GetNewsByWeek() ([]d.News, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var weekNews []d.News

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1) // Monday
	if now.Weekday() == time.Sunday {
		weekStart = now.AddDate(0, 0, -6) // Special case for Sunday
	}
	monday := weekStart.Truncate(24 * time.Hour)
	fmt.Printf("monday: %v\n", monday)

	findOptions := options.Find().
		SetProjection(bson.M{"_id": 0, "tag": 0, "status": 0, "content_type": 0, "updated_at": 0}).
		SetLimit(4).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	filter := bson.M{
		"status":         true,
		"content_status": "published",
		"content_type":   "general",
		"created_at":     bson.M{"$gte": monday},
	}

	cursor, err := n.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &weekNews); err != nil {
		log.Printf("Error decoding repo last news: %v", err)
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return weekNews, nil

}

// /// get area ///

func (n *MongoNewsRepository) UpdateNews(id string, news *d.News) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"title":          news.Title,
			"abstract":       news.Abstract,
			"detail":         news.Detail,
			"image":          news.Image,
			"category_id":    news.CategoryID,
			"tag":            news.Tag,
			"status":         news.Status,
			"content_status": news.ContentStatus,
			"content_type":   news.ContentType,
			"updated_at":     news.UpdatedAt,
		},
	}

	_, err = n.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		update,
		options.Update().SetUpsert(true),
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

func (n *MongoNewsRepository) DeleteImg(path string) error {
	fullPath := "./upload/news_image/" + path

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("Mongo : Image not found at: %s", fullPath)
		return nil
	}

	err := os.Remove(fullPath)
	if err != nil {
		log.Printf(" Failed to delete image at %s: %v", fullPath, err)
		return err
	}

	log.Printf("Image deleted: %s", fullPath)
	return nil
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
