// secondary adapters //
package repository

import (
	d "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoNewsRepository struct {
	db *mongo.Collection
}

func NewNewsRepo(db *mongo.Database) port.NewsRepository {
	return &MongoNewsRepository{db: db.Collection("news")}
}

func (n *MongoNewsRepository) Create(news *d.News) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	news.ID = primitive.NewObjectID()
	_, err := n.db.InsertOne(ctx, news)

	return err
}

func (n *MongoNewsRepository) GetNewsPagination(page int, limit int) ([]d.News, int, error) {

	var news []d.News
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// คำนวณการข้ามข้อมูลระหว่างหน้า ข้ามหน้า //
	skip := (page - 1) * limit

	findOption := options.Find()
	findOption.SetSkip(int64(skip))
	findOption.SetLimit(int64(limit))

	cursor, err := n.db.Find(ctx, bson.M{}, findOption)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var new d.News
		if err := cursor.Decode(&new); err != nil {
			return nil, 0, err
		}
		news = append(news, new)
	}

	total, err := n.db.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return news, int(total), nil
}

func (n *MongoNewsRepository) GetNewsByID(id string) (*d.News, error) {
	var news d.News
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// แปลงค่า string ที่อยู่ในรูปแบบ Hexadecimal ให้กลายเป็น primitive.ObjectID //
	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// find id in database //
	err = n.db.FindOne(ctx, bson.M{"_id": ObjID}).Decode(&news)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &news, nil
}

func (n *MongoNewsRepository) UpdateNews(id string, news *d.News) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"title":       news.Title,
			"detail":      news.Detail,
			"image":       news.Image,
			"category_id": news.CategoryID,
			"tag":         news.Tag,
			"updated_at":  news.UpdatedAt,
		},
	}

	_, err = n.db.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		update,
		options.Update().SetUpsert(true),
	)

	return err
}

func (n *MongoNewsRepository) Delete(id string) error {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result := n.db.FindOneAndDelete(ctx, bson.M{"_id": objID})
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
