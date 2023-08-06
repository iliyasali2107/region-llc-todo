package db

import (
	"context"
	"log"

	"region-llc-todo/pkg/config"
	"region-llc-todo/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	InsertTodo(ctx context.Context, todo models.Todo) error
}

type storage struct {
	DB             *mongo.Client
	TodoCollection *mongo.Collection
}

// TODO: close connection

func Init(cfg config.Config) Storage {
	ctx := context.Background()

	opts := options.Client().ApplyURI(cfg.DBUrl)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
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
func (s *storage) InsertTodo(ctx context.Context, todo models.Todo) error {
	_, err := s.TodoCollection.InsertOne(ctx, todo)
	if err != nil {
		return err
	}

	return nil
}
