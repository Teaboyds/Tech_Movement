package mongodb

import (
	"backend_tech_movement_hex/internal/adapter/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	MongoClient *mongo.Client
	DBName      string
}

func ConnectDB(ctx context.Context, config *config.DB) (*Database, error) {

	mongoURL := config.URL

	clientOption := options.Client().ApplyURI(mongoURL)
	// .SetAuth(options.Credential{
	// 		AuthSource:    os.Getenv("DB_AUTHSOURCE"),
	// 		Username:      os.Getenv("DB_USERNAME"),
	// 		Password:      os.Getenv("MONGO_PASSWORD"),
	// 		AuthMechanism: os.Getenv("DB_AUTHMEC"),
	// 	})
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Can't connect mongodb: ", err)
	} else {
		log.Println("Connected to Mongodb")
	}

	return &Database{MongoClient: client, DBName: config.DB_NAME}, nil
}

func (db *Database) Close(ctx context.Context) error {
	if err := db.MongoClient.Disconnect(ctx); err != nil {
		log.Println("Database close error -:", err)
		return err
	}
	log.Println("Database close!")
	return nil
}

func (db *Database) Collection(name string) *mongo.Collection {
	return db.MongoClient.Database(db.DBName).Collection(name)
}
