package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func LoadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Env coudln't Load", err)
	}
}

func ConnectDB() {

	LoadEnv()

	mongoURL := os.Getenv("MONGO_URI")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOption := options.Client().ApplyURI(mongoURL)
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

	MongoClient = client
}

func GetDatabase() *mongo.Database {
	return MongoClient.Database("hexagonal_db")
}
