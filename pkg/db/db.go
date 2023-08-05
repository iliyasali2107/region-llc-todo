package db

import (
	"context"
	"fmt"
	"log"

	"region-llc-todo/pkg/config"
	"region-llc-todo/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	InsertTodo(ctx context.Context, todo models.Todo) (int64, error)
}

type storage struct {
	DB             *mongo.Client
	TodoCollection *mongo.Collection
}

// TODO: close connection

func Init(cfg config.Config) Storage {
	ctx := context.Background()

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.DBUrl).
		SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(ctx, opts)
	fmt.Println("qewrqwer")
	if err != nil {
		log.Fatal(err)
	}

	todoCollection := InitCollection(client, cfg)
	return &storage{
		DB:             client,
		TodoCollection: todoCollection,
	}
}

func InitCollection(db *mongo.Client, cfg config.Config) *mongo.Collection {
	ctx := context.Background()
	todoCollection := db.Database(cfg.DbName).Collection(cfg.TodoCollectionName)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: 1},
			{Key: "active_at", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := todoCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return db.Database(cfg.DbName).Collection(cfg.TodoCollectionName)
}

// TODO: duplicate field error
func (s *storage) InsertTodo(ctx context.Context, todo models.Todo) (int64, error) {
	insertResult, err := s.TodoCollection.InsertOne(ctx, todo)
	if err != nil {
		return 0, err
	}

	n, ok := insertResult.InsertedID.(int64)
	if !ok {
		return 0, fmt.Errorf("couldn't assert type int64")
	}

	return n, nil
}
