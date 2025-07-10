package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var CronjobsCollection *mongo.Collection

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:12345abc@localhost:27017/?directConnection=true&authMechanism=SCRAM-SHA-1&authSource=admin"))
	if err != nil {
		log.Fatal("❌ Error connecting to MongoDB:", err)
	}

	CronjobsCollection = Client.Database("cronjobdb").Collection("jobs")
	log.Println("✅ Connected to MongoDB")
}
