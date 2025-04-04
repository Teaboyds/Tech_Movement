package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	dt "backend_tech_movement_hex/internal/core/domain"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTagsRepository struct {
	collection *mongo.Collection
}

func NewTagRepo(db mongodb.Database) *MongoTagsRepository {
	return &MongoTagsRepository{collection: db.Collection("tags")}
}

func (t *MongoTagsRepository) SavaTags(tags dt.Tags) error {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tags.ID = primitive.NewObjectID()

	_, err := t.collection.InsertOne(ctx, tags)
	if err != nil {
		log.Printf("Failed to insert tag: %v", err)
		return err
	}

	return nil
}

func (t *MongoTagsRepository) GetTagsById(id string) (*dt.Tags, error) {

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var tags dt.Tags

	filter := bson.M{"_id": ObjID}

	err = t.collection.FindOne(ctx, filter).Decode(&tags)
	if err != nil {
		return nil, err
	}

	return &tags, nil
}
