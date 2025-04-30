package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCategoryRepository struct {
	db *mongo.Collection
}

func NewCategoryRepositoryMongo(db *mongodb.Database) port.CategoryRepository {
	return &MongoCategoryRepository{db: db.Collection("categories")}
}

func (cat *MongoCategoryRepository) Create(category *domain.CategoryRequest) error {

	cateDoc := &models.MongoCategory{
		Name:         category.Name,
		CategoryType: category.CategoryType,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := cat.db.InsertOne(context.Background(), cateDoc)
	return err
}

func (cat *MongoCategoryRepository) GetByID(id string) (*domain.CategoryResponse, error) {

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var category models.MongoCategory

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	filter := bson.M{"_id": ObjID}

	err = cat.db.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	response := &domain.CategoryResponse{
		ID:   category.ID.Hex(),
		Name: category.Name,
	}

	return response, nil
}

func (cat *MongoCategoryRepository) GetByIDs(ids []string) ([]*domain.CategoryResponse, error) {

	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs = append(objIDs, objID)
	}

	var categories []models.MongoCategory

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	filter := bson.M{
		"_id": bson.M{"$in": objIDs},
	}

	cursor, err := cat.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	var responses []*domain.CategoryResponse
	for _, category := range categories {
		resp := &domain.CategoryResponse{
			ID:   category.ID.Hex(),
			Name: category.Name,
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

func (cat *MongoCategoryRepository) GetByName(name string) (*domain.Category, error) {

	var category models.MongoCategory

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	err := cat.db.FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("category with name %s not found", name)
		}
		return nil, fmt.Errorf("failed to find category by name: %w", err)
	}

	resp := &domain.Category{
		ID:   category.ID.Hex(),
		Name: category.Name,
	}

	return resp, nil
}

func (cat *MongoCategoryRepository) GetAll() ([]domain.CategoryResponse, error) {
	var categories []models.MongoCategory
	cursor, err := cat.db.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var category models.MongoCategory
		if err := cursor.Decode(&category); err != nil {
			log.Println("Error decoding category:", err)
			continue
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	response := make([]domain.CategoryResponse, len(categories))
	for i, c := range categories {
		response[i] = domain.CategoryResponse{
			ID:   c.ID.Hex(),
			Name: c.Name,
		}
	}

	fmt.Printf("response: %v\n", response)

	return response, nil
}

func (cat *MongoCategoryRepository) UpdateCategory(id string, category *domain.Category) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":       category.Name,
			"updated_at": time.Now(),
		},
	}

	_, err = cat.db.UpdateOne(
		ctx,
		bson.M{"_id": objId},
		update,
		options.Update().SetUpsert(false),
	)

	return err
}

func (cat *MongoCategoryRepository) DeleteCategory(id string) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result := cat.db.FindOneAndDelete(ctx, bson.M{"_id": ObjId})
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (cat *MongoCategoryRepository) ExistsByName(name string) (bool, error) {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	// ค้นหา database ว่ามี name ซ้ำกันยุบ่ //
	filter := bson.M{"name": name}
	var result bson.M

	err := cat.db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
