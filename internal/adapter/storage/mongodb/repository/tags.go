package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	dt "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/utils"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTagsRepository struct {
	collection *mongo.Collection
}

func NewTagRepo(db mongodb.Database) *MongoTagsRepository {
	return &MongoTagsRepository{collection: db.Collection("tags")}
}

func (t *MongoTagsRepository) SavaTags(tags dt.Tags) error {

	ctx, cancel := utils.NewTimeoutContext()
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

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var tags dt.Tags

	filter := bson.M{"_id": ObjID}

	err = t.collection.FindOne(ctx, filter).Decode(&tags)
	if err != nil {
		return nil, err
	}

	return &tags, nil
}

func (t *MongoTagsRepository) GetAllTags() ([]dt.Tags, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var tags []dt.Tags
	cursor, err := t.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, bson.ErrDecodeToNil
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var loopTags dt.Tags
		if err := cursor.Decode(&loopTags); err != nil {
			log.Println("Error decoding category:", err)
			continue
		}
		tags = append(tags, loopTags)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (t *MongoTagsRepository) EditTags(id string, tags dt.Tags) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name": tags.Name,
		},
	}

	_, err = t.collection.UpdateOne(
		ctx,
		bson.M{"_id": ObjId},
		update,
		options.Update().SetUpsert(true),
	)

	return err
}

func (t *MongoTagsRepository) DeleteTags(id string) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result := t.collection.FindOneAndDelete(ctx, bson.M{"_id": ObjId})
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (t *MongoTagsRepository) ExistsByName(name string) (bool, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	// ค้นหา database ว่ามี name ซ้ำกันยุบ่ //
	filter := bson.M{"name": name}
	var result bson.M

	err := t.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
