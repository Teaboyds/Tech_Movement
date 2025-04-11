package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/core/domain"
	dt "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTagsRepository struct {
	collection *mongo.Collection
}

func NewTagRepo(db *mongodb.Database) port.TagsRepository {
	return &MongoTagsRepository{collection: db.Collection("tags")}
}

func (t *MongoTagsRepository) SavaTags(tags *dt.Tags) error {

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

// look up get ตอน ข่าวใหย่ //
func (t *MongoTagsRepository) GetTagsByIdArray(id []string) ([]*dt.Tags, error) {

	var ObjIDs []primitive.ObjectID

	// loop id แล้วเก็บว้ายนาย ObjIDs //
	for _, ids := range id {
		objID, err := primitive.ObjectIDFromHex(ids)
		if err != nil {
			return nil, err
		}
		ObjIDs = append(ObjIDs, objID)
	}

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	// filter หา id ใน ObjIDs //
	filter := bson.M{"_id": bson.M{"$in": ObjIDs}}

	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var tags []*dt.Tags
	for cursor.Next(ctx) {
		var tag dt.Tags
		if err := cursor.Decode(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tags, err
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

func (t *MongoTagsRepository) EditTags(id string, tags *dt.Tags) error {

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

func (t *MongoTagsRepository) GetByName(name string) ([]dt.Tags, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var tags []domain.Tags
	err := t.collection.FindOne(ctx, bson.M{"name": name}).Decode(&tags)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return tags, err
}

func (t *MongoTagsRepository) GetByID(id string) (*dt.Tags, error) {

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID%w", err)
	}

	var tags domain.Tags

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	filter := bson.M{"_id": ObjID}

	err = t.collection.FindOne(ctx, filter).Decode(&tags)
	if err != nil {
		return nil, err
	}

	return &tags, nil
}
