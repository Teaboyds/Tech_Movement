package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	mongoUtils "backend_tech_movement_hex/internal/adapter/storage/mongodb/utils"
	"backend_tech_movement_hex/internal/core/domain"
	up "backend_tech_movement_hex/internal/core/domain"
	"backend_tech_movement_hex/internal/core/port"
	"backend_tech_movement_hex/internal/core/utils"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUploadRepository struct {
	collection *mongo.Collection
}

func NewUploadRepo(db *mongodb.Database) port.UploadRepository {
	return &MongoUploadRepository{collection: db.Collection("file_Upload")}
}

func (n *MongoUploadRepository) EnsureFileIndexs() error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "created_at", Value: -1}, // ใหม่ไปเก่า //
			{Key: "_id", Value: -1},
		},
	}

	_, err := n.collection.Indexes().CreateOne(ctx, indexModel)

	return err
}

func (ul *MongoUploadRepository) SaveImage(file *up.UploadFile) error {

	ctx, cancle := utils.NewTimeoutContext()
	defer cancle()

	responseFile := &models.MongoUploadRepository{
		Path:      file.Path,
		Name:      file.Name,
		FileType:  file.FileType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := ul.collection.InsertOne(ctx, responseFile)

	return err
}

func (ul *MongoUploadRepository) GetFileByID(id string) (*up.UploadFile, error) {

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var file models.MongoUploadRepository

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	filter := bson.M{"_id": ObjID}

	err = ul.collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		return nil, err
	}

	response := &domain.UploadFile{
		ID:       file.ID.Hex(),
		Path:     file.Path,
		Name:     file.Name,
		FileType: file.FileType,
	}

	return response, nil
}

func (ul *MongoUploadRepository) GetFilesByIDs(ids []string) ([]up.UploadFile, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var objIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid image ID format: %w", err)
		}
		objIDs = append(objIDs, objID)
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}

	cursor, err := ul.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding files: %w", err)
	}
	defer cursor.Close(ctx)

	var files []up.UploadFile
	for cursor.Next(ctx) {
		var file models.MongoUploadRepository
		if err := cursor.Decode(&file); err != nil {
			return nil, fmt.Errorf("error decoding file: %w", err)
		}
		files = append(files, up.UploadFile{
			ID:       file.ID.Hex(),
			Path:     file.Path,
			Name:     file.Name,
			FileType: file.FileType,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	if len(files) != len(ids) {
		return nil, fmt.Errorf("some image IDs not found in database")
	}

	return files, nil
}

func (ul *MongoUploadRepository) GetAllFile() ([]up.UploadFile, error) {
	var files []models.MongoUploadRepository
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	filter := bson.M{}
	findOptions := options.Find().
		SetLimit(10).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := ul.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var file models.MongoUploadRepository
		if err := cursor.Decode(&file); err != nil {
			log.Println("Error decoding category:", err)
			continue
		}
		files = append(files, file)
	}
	defer cursor.Close(ctx)

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	response := make([]domain.UploadFile, len(files))
	for i, c := range files {
		response[i] = domain.UploadFile{
			ID:        c.ID.Hex(),
			Path:      c.Path,
			Name:      c.Name,
			FileType:  c.FileType,
			CreatedAt: c.CreatedAt.In(time.FixedZone("UTC+7", 7*60*60)).Format("2006-01-02 15:04:05"),
			UpdatedAt: c.UpdatedAt.In(time.FixedZone("UTC+7", 7*60*60)).Format("2006-01-02 15:04:05"),
		}
	}

	return response, err

}

// check if not me obj kub //
func (ul *MongoUploadRepository) ValidateImageIDs(ids []string) ([]string, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var objIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid image ID format: %w", err)
		}
		objIDs = append(objIDs, objID)
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	count, err := ul.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error checking image IDs: %w", err)
	}

	if count != int64(len(ids)) {
		return nil, fmt.Errorf("some image IDs not found")
	}

	return ids, nil
}

func (ul *MongoUploadRepository) DeleteFile(id string) error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result := ul.collection.FindOneAndDelete(ctx, bson.M{"_id": ObjId})
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (ul *MongoUploadRepository) GetFilesByIDsVTest(ids []string) ([]*up.UploadFile, error) {

	obj, err := mongoUtils.ConvertStringToObjectIDArray(ids)
	if err != nil {
		return nil, err
	}

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": obj}}

	cursor, err := ul.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding files: %w", err)
	}
	defer cursor.Close(ctx)

	var files []*up.UploadFile
	for cursor.Next(ctx) {
		var file models.MongoUploadRepository
		if err := cursor.Decode(&file); err != nil {
			return nil, fmt.Errorf("error decoding file: %w", err)
		}
		files = append(files, &up.UploadFile{
			ID:        file.ID.Hex(),
			Path:      file.Path,
			Name:      file.Name,
			FileType:  file.FileType,
			CreatedAt: file.CreatedAt.Format(time.RFC3339),
			UpdatedAt: file.UpdatedAt.Format(time.RFC3339),
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return files, nil
}
