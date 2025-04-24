package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	dm "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMediaRepository struct {
	collection *mongo.Collection
}

func NewMediaRepo(db *mongodb.Database) port.MediaRepository {
	return &MongoMediaRepository{collection: db.Collection("media")}
}

func (n *MongoMediaRepository) EnsureMediaIndexs() error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "created_at", Value: -1}, // ใหม่ไปเก่า //
			{Key: "category_id._id", Value: -1},
			{Key: "_id", Value: -1},
		},
	}

	_, err := n.collection.Indexes().CreateOne(ctx, indexModel)

	return err
}

func (m *MongoMediaRepository) SaveMedia(media *dm.Media) error {

	ctx, cancle := utils.NewTimeoutContext()
	defer cancle()

	media.ID = primitive.NewObjectID()
	_, err := m.collection.InsertOne(ctx, media)
	return err
}

func (m *MongoMediaRepository) GetLastMedia() ([]dm.Media, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var lastMedia []dm.Media

	findOptions := options.Find().
		SetProjection(bson.M{"_id": 0, "tag": 0, "status": 0, "detail": 0, "updated_at": 0}).
		SetLimit(4).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	fileter := bson.M{}

	cursor, err := m.collection.Find(ctx, fileter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &lastMedia); err != nil {
		log.Printf("Error decoding repo last news: %v", err)
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("lastMedia: %v\n", lastMedia)
	return lastMedia, nil
}
