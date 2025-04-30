// secondary adapters //
package repository

import (
	"backend_tech_movement_hex/internal/adapter/storage/mongodb"
	"backend_tech_movement_hex/internal/adapter/storage/mongodb/models"
	"backend_tech_movement_hex/internal/core/domain"
	d "backend_tech_movement_hex/internal/core/domain"
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

type MongoNewsRepository struct {
	collection   *mongo.Collection
	categoryRepo port.CategoryRepository
	uploadRepo   port.UploadRepository
}

func NewNewsRepo(db *mongodb.Database, categoryRepo port.CategoryRepository, uploadRepo port.UploadRepository) port.NewsRepository {
	return &MongoNewsRepository{
		collection:   db.Collection("news"),
		categoryRepo: categoryRepo,
		uploadRepo:   uploadRepo,
	}
}

// // index model ////
func (n *MongoNewsRepository) EnsureNewsIndexs() error {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "_id", Value: -1},
		},
		Options: options.Index().SetName("status_id_desc"),
	}

	_, err := n.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

// /// create area ///
func (n *MongoNewsRepository) SaveNews(news *d.NewsRequest) error {

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	Category, err := n.categoryRepo.GetByID(news.Category)
	if err != nil {
		return err
	}

	CateOBJ, err := primitive.ObjectIDFromHex(Category.ID)
	if err != nil {
		return err
	}

	File, err := n.uploadRepo.ValidateImageIDs(news.Image)
	if err != nil {
		return err
	}

	var imageObjectIDs []primitive.ObjectID
	for _, id := range File {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("error converting image ID: %w", err)
		}
		imageObjectIDs = append(imageObjectIDs, objID)
	}

	newDoc := &models.MongoNews{
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		Image:       imageObjectIDs,
		CategoryID:  CateOBJ,
		Tag:         news.Tag,
		Status:      news.Status,
		ContentType: news.ContentType,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = n.collection.InsertOne(ctx, newDoc)
	return err
}

// /// get area ///

// func (n *MongoNewsRepository) GetNewsPagination(lastID string, limit int) ([]d.News, error) {

// 	var news []d.News
// 	ctx, cancel := utils.NewTimeoutContext()
// 	defer cancel()

// 	// cursor based pagination //
// 	filter := bson.M{}
// 	if lastID != "" {
// 		ObjID, err := primitive.ObjectIDFromHex(lastID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		filter["_id"] = bson.M{"$lt": ObjID}
// 	}

// 	// set spect for sort //
// 	findOption := options.Find()
// 	findOption.SetSort(bson.M{"_id": -1})
// 	findOption.SetLimit(int64(limit))

// 	cursor, err := n.collection.Find(ctx, filter, findOption)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// loop decode and appened//
// 	for cursor.Next(ctx) {
// 		var new d.News
// 		if err := cursor.Decode(&new); err != nil {
// 			return nil, err
// 		}
// 		news = append(news, new)
// 	}
// 	defer cursor.Close(ctx)

// 	return news, nil
// }

func (n *MongoNewsRepository) GetNewsByID(id string) (*d.NewsResponse, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var newsResp models.MongoNews
	err = n.collection.FindOne(ctx, bson.M{"_id": ObjID}).Decode(&newsResp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	var images []domain.UploadFileResponse
	if len(newsResp.Image) > 0 {
		imgIDs := make([]string, len(newsResp.Image))
		for i, id := range newsResp.Image {
			imgIDs[i] = id.Hex()
		}

		uploadFiles, err := n.uploadRepo.GetFilesByIDs(imgIDs)
		if err != nil {
			return nil, err
		}

		for _, f := range uploadFiles {
			images = append(images, domain.UploadFileResponse{
				ID:       f.ID,
				Path:     f.Path,
				Name:     f.Name,
				FileType: f.FileType,
			})
		}
	}

	var categoryResponse domain.CategoryResponse
	if newsResp.CategoryID != primitive.NilObjectID {
		category, err := n.categoryRepo.GetByID(newsResp.CategoryID.Hex())
		if err != nil {
			return nil, err
		}
		categoryResponse = *category
	}

	response := &domain.NewsResponse{
		ID:          newsResp.ID.Hex(),
		Title:       newsResp.Title,
		Description: newsResp.Description,
		Content:     newsResp.Content,
		Image:       images,
		CategoryID:  categoryResponse,
		Tag:         newsResp.Tag,
		Status:      newsResp.Status,
		ContentType: newsResp.ContentType,
		CreatedAt:   newsResp.CreatedAt.Format(time.RFC1123),
		UpdatedAt:   newsResp.UpdatedAt.Format(time.RFC1123),
	}

	fmt.Printf("response: %v\n", response)

	return response, nil
}

func (n *MongoNewsRepository) GetLastNews() ([]*domain.HomePageLastedNewResponse, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	var lastNews []models.MongoNews

	findOptions := options.Find().
		SetProjection(bson.M{"_id": 0, "tag": 0, "status": 0, "updated_at": 0}).
		SetLimit(5).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	filter := bson.M{
		"status": true,
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

	imageIDMap := make(map[string]struct{})
	categoryIDMap := make(map[string]struct{})

	for _, news := range lastNews {
		for _, imgID := range news.Image {
			imageIDMap[imgID.Hex()] = struct{}{}
		}
		if news.CategoryID != primitive.NilObjectID {
			categoryIDMap[news.CategoryID.Hex()] = struct{}{}
		}
	}

	imageIDs := make([]string, 0, len(imageIDMap))
	for id := range imageIDMap {
		imageIDs = append(imageIDs, id)
	}

	categoryIDs := make([]string, 0, len(categoryIDMap))
	for id := range categoryIDMap {
		categoryIDs = append(categoryIDs, id)
	}

	uploadFiles, err := n.uploadRepo.GetFilesByIDs(imageIDs)
	if err != nil {
		return nil, err
	}

	categories, err := n.categoryRepo.GetByIDs(categoryIDs)
	if err != nil {
		return nil, err
	}

	uploadFileMap := make(map[string]domain.UploadFileResponseHomePage)
	for _, f := range uploadFiles {
		uploadFileMap[f.ID] = domain.UploadFileResponseHomePage{
			ID:       f.ID,
			Path:     f.Path,
			Filetype: f.FileType,
		}
	}

	categoryMap := make(map[string]domain.CategoryResponse)
	for _, ca := range categories {
		categoryMap[ca.ID] = *ca
	}

	var responseNews []*domain.HomePageLastedNewResponse

	for _, news := range lastNews {
		var images []domain.UploadFileResponseHomePage
		for _, imgID := range news.Image {
			if img, ok := uploadFileMap[imgID.Hex()]; ok {
				images = append(images, img)
			}
		}

		var categoryResponse domain.CategoryResponse
		if news.CategoryID != primitive.NilObjectID {
			if cat, ok := categoryMap[news.CategoryID.Hex()]; ok {
				categoryResponse = cat
			}
		}

		resp := &domain.HomePageLastedNewResponse{
			Title:       news.Title,
			Detail:      news.Description,
			Image:       images,
			Category:    categoryResponse,
			ContentType: news.ContentType,
			CreatedAt:   news.CreatedAt,
		}
		responseNews = append(responseNews, resp)
	}

	return responseNews, nil
}

func (n *MongoNewsRepository) GetTechnologyNews() ([]*d.HomePageLastedNewResponse, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "status", Value: true},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "category_id", Value: 1},
			{Key: "created_at", Value: -1},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category_id"},
			{Key: "news", Value: bson.D{
				{Key: "$first", Value: "$$ROOT"},
			}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "categories"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "category"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$category"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := n.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		News     models.MongoNews     `bson:"news"`
		Category models.MongoCategory `bson:"category"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var responses []*d.HomePageLastedNewResponse

	for _, result := range results {
		var images []d.UploadFileResponseHomePage
		for _, imgID := range result.News.Image {
			file, err := n.uploadRepo.GetFileByID(imgID.Hex())
			if err == nil {
				images = append(images, d.UploadFileResponseHomePage{
					ID:       file.ID,
					Path:     file.Path,
					Filetype: file.FileType,
				})
			}
		}

		categoryResp := d.CategoryResponse{
			ID:   result.Category.ID.Hex(),
			Name: result.Category.Name,
		}

		responses = append(responses, &d.HomePageLastedNewResponse{
			Title:       result.News.Title,
			Detail:      result.News.Description,
			Image:       images,
			Category:    categoryResp,
			ContentType: result.News.ContentType,
			CreatedAt:   result.News.CreatedAt,
		})
	}

	return responses, nil
}

// func (n *MongoNewsRepository) GetNewsByCategoryHomePage(categoryID string) ([]d.News, error) {

// 	ctx, cancel := utils.NewTimeoutContext()
// 	defer cancel()

// 	var newsCategory []d.News

// 	filter := bson.M{
// 		"status":         true,
// 		"content_status": "published",
// 		"content_type":   "general",
// 	}

// 	if categoryID != "" {
// 		objID, err := primitive.ObjectIDFromHex(categoryID)
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid CategoryID: %v", err)
// 		}
// 		filter["category_id._id"] = objID
// 	}

// 	findOptions := options.Find().
// 		SetProjection(bson.M{
// 			"_id":          0,
// 			"tag":          0,
// 			"status":       0,
// 			"content_type": 0,
// 			"updated_at":   0,
// 		}).
// 		SetLimit(9).
// 		SetSort(bson.D{{Key: "_id", Value: -1}})

// 	cursor, err := n.collection.Find(ctx, filter, findOptions)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot fetching news: %v", err)
// 	}
// 	if cursor != nil {
// 		defer cursor.Close(ctx)
// 	}

// 	if err := cursor.All(ctx, &newsCategory); err != nil {
// 		return nil, fmt.Errorf("error decoding news: %v", err)
// 	}

// 	return newsCategory, nil
// }

// /// get area ///

// func (n *MongoNewsRepository) UpdateNews(id string, news *d.News) error {
// 	ctx, cancel := utils.NewTimeoutContext()
// 	defer cancel()

// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"title":        news.Title,
// 			"description":  news.Description,
// 			"content":      news.Content,
// 			"image":        news.Image,
// 			"category_id":  news.CategoryID,
// 			"tag":          news.Tag,
// 			"status":       news.Status,
// 			"content_type": news.ContentType,
// 			"updated_at":   news.UpdatedAt,
// 		},
// 	}

// 	_, err = n.collection.UpdateOne(
// 		ctx,
// 		bson.M{"_id": objID},
// 		update,
// 		options.Update().SetUpsert(true),
// 	)

// 	return err
// }

// func (n *MongoNewsRepository) Delete(id string) error {
// 	ctx, cancle := utils.NewTimeoutContext()
// 	defer cancle()

// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := n.collection.DeleteOne(ctx, bson.M{"_id": objID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		log.Printf("No news found with ID: %s", id)
// 		return errors.New("news not found or already deleted")
// 	}

// 	return nil
// }

// func (n *MongoNewsRepository) DeleteImg(path string) error {
// 	fullPath := "./upload/news_image/" + path

// 	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
// 		log.Printf("Mongo : Image not found at: %s", fullPath)
// 		return nil
// 	}

// 	err := os.Remove(fullPath)
// 	if err != nil {
// 		log.Printf(" Failed to delete image at %s: %v", fullPath, err)
// 		return err
// 	}

// 	log.Printf("Image deleted: %s", fullPath)
// 	return nil
// }

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
