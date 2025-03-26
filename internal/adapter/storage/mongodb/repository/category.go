package repository

import (
	"backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCategoryRepository struct {
	db *mongo.Collection
}

func NewCategoryRepositoryMongo(db *mongo.Database) port.CategoryRepository {
	return &MongoCategoryRepository{db: db.Collection("categories")}
}

func (cat *MongoCategoryRepository) Create(category *domain.Category) error {

	category.ID = primitive.NewObjectID()
	_, err := cat.db.InsertOne(context.Background(), category)
	return err
}

func (cat *MongoCategoryRepository) GetByID(id string) (*domain.Category, error) {

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var category domain.Category

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.M{"_id": ObjID}

	err = cat.db.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (cat *MongoCategoryRepository) GetByName(name string) (*domain.Category, error) {
	var category domain.Category
	err := cat.db.FindOne(context.Background(), bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}

		return nil, err
	}

	return &category, nil
}

func (cat *MongoCategoryRepository) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	cursor, err := cat.db.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var category domain.Category
		if err := cursor.Decode(&category); err != nil {
			log.Println("Error decoding category:", err)
			continue
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
