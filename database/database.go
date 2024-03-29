package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ad/model"
	"os"
)

type MongoDB struct {
	DB                    *mongo.Client
	AdCollections         *mongo.Collection
	CurrentAdsCollections *mongo.Collection
}

func New() (*mongo.Client, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	uri := "mongodb://" + user + ":" + pass + "@" + host + ":" + port

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	coll := client.Database("dcard_ads").Collection("current_ads")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "conditions.country", Value: 1},
			{Key: "conditions.gender", Value: 1},
			{Key: "conditions.platform", Value: 1},
			{Key: "conditions.age", Value: 1},
		},
	}
	name, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic(err)
	}
	log.Println("Name of Index Created: " + name)

	return client, err
}

func NewTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Ad{})

	return db, nil
}
