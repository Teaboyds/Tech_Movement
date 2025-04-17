// secondary adapters //
package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/utils"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoNewsRepository struct {
	collection *mongo.Collection
	redis      *redis.Client
}

func NewNewsRepo(db *mongodb.Database, redisClient *redis.Client) *MongoNewsRepository {
	return &MongoNewsRepository{
		collection: db.Collection("news"),
		redis:      redisClient,
	}
}

func (n *MongoNewsRepository) Create(news *d.News) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	news.ID = primitive.NewObjectID()
	_, err := n.collection.InsertOne(ctx, news)

	return err
}

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

	// แปลงค่า string ที่อยู่ในรูปแบบ Hexadecimal ให้กลายเป็น primitive.ObjectID //
	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// debug if not found database or cache server der kub //
	if n.collection == nil {
		return nil, errors.New("mongoDB client is nil")
	}
	if n.redis == nil {
		return nil, errors.New("redis client is nil")
	}

	cacheKey := "News_Keys_" + id
	val, err := n.redis.Get(ctx, cacheKey).Result()

	// if บ่เจอ cache เด้อครับ //
	if err == redis.Nil {
		// find id in database //
		err = n.collection.FindOne(ctx, bson.M{"_id": ObjID}).Decode(&news)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, err
			}
			return nil, err
		}

		// แปลงร่าง struct เป็น json //
		data, err := json.Marshal(news)
		if err != nil {
			return nil, err
		}

		// ละกะเซ็ทแคชไว้บาดนิ สิบนาที //
		err = n.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// ถ้าพ้อ cache กะแปลง jsonมาเป็น struct
		err = json.Unmarshal([]byte(val), &news)
		if err != nil {
			log.Println("Error unmarshaling data from Redis:", err)
			return nil, err
		}
	}

	return &news, nil
}

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
			"detail":         news.Detail,
			"image":          news.Image,
			"category_id":    news.CategoryID,
			"tag":            news.Tag,
			"status":         news.Status,
			"content_status": news.ContentStatus,
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
	fullPath := "./upload/image/" + path // หรือใช้ absolute root path หากต้องการ

	// เช็คก่อนว่าไฟล์มีอยู่หรือไม่
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("⚠️ Image not found at: %s", fullPath)
		return nil
	}

	// ลบไฟล์
	err := os.Remove(fullPath)
	if err != nil {
		log.Printf("❌ Failed to delete image at %s: %v", fullPath, err)
		return err
	}

	log.Printf("✅ Image deleted: %s", fullPath)
	return nil
}

func (n *MongoNewsRepository) GetNewsByCategory(CategoryId string) ([]d.News, error) {

	var newsList []d.News
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjID, err := primitive.ObjectIDFromHex(CategoryId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"category_id._id": ObjID,
	}
	cursor, err := n.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &newsList); err != nil {
		return nil, err
	}

	return newsList, nil
}

func (n *MongoCategoryRepository) GetNewsByTags(name string) ([]d.News, error) {

	var news []d.News
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	filter := bson.M{"tags": name}
	cursor, err := n.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &news); err != nil {
		return nil, err
	}

	return news, nil
}
