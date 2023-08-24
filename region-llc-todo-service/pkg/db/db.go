package db

import (
	"context"
	"log"
	"time"

	"region-llc-todo-service/pkg/config"
	"region-llc-todo-service/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	InsertTodo(ctx context.Context, todo models.Todo) (string, error)
	UpdateTodoById(ctx context.Context, todo models.Todo) (int64, error)
	DeleteTodoById(ctx context.Context, id string) (int64, error)
	UpdateAsDone(ctx context.Context, id string) (int64, error)
	GetTodosByFilterDone(ctx context.Context) ([]models.Todo, error)
	GetTodosByFilterActive(ctx context.Context) ([]models.Todo, error)
	GetOneTodo(ctx context.Context, id string) (todo models.Todo, err error)
	Close(ctx context.Context) error
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

func (s *storage) Close(ctx context.Context) error {
	return s.DB.Disconnect(context.Background())
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

func (s *storage) InsertTodo(ctx context.Context, todo models.Todo) (string, error) {
	res, err := s.TodoCollection.InsertOne(ctx, todo)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrDuplicate
		}

		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *storage) UpdateTodoById(ctx context.Context, todo models.Todo) (int64, error) {
	objId, _ := primitive.ObjectIDFromHex(todo.Id)
	filter := bson.M{"_id": objId}
	update := bson.M{"$set": bson.M{"title": todo.Title, "active_at": todo.ActiveAt}}
	res, err := s.TodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, ErrNotFound
		}
		return 0, err
	}

	if res.MatchedCount == 0 {
		return 0, ErrNotFound
	}

	if res.ModifiedCount == 0 {
		return 0, ErrModify
	}

	return res.ModifiedCount, nil
}

func (s *storage) DeleteTodoById(ctx context.Context, id string) (int64, error) {
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	res, err := s.TodoCollection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	if res.DeletedCount == 0 {
		return 0, ErrNotFound
	}

	return res.DeletedCount, nil
}

func (s *storage) UpdateAsDone(ctx context.Context, id string) (int64, error) {
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	update := bson.M{"$set": bson.M{"status": "done"}}
	res, err := s.TodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, ErrNotFound
		}

		return 0, err
	}

	if res.MatchedCount == 0 {
		return 0, ErrNotFound
	}

	if res.ModifiedCount == 0 {
		return 0, ErrModify
	}

	return res.ModifiedCount, nil
}

func (s *storage) GetTodosByFilterActive(ctx context.Context) ([]models.Todo, error) {
	sortOptions := options.Find().SetSort(bson.M{"active_at": 1})
	today := time.Now().UTC()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	filter := bson.M{"status": StatusActive, "active_at": bson.M{"$lte": today}}

	cursor, err := s.TodoCollection.Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, err
	}

	var todos []models.Todo
	err = cursor.All(ctx, &todos)
	if err != nil {
		return nil, ErrNotFound
	}

	return todos, nil
}

func (s *storage) GetTodosByFilterDone(ctx context.Context) ([]models.Todo, error) {
	filter := bson.M{"status": StatusDone}

	sortOptions := options.Find().SetSort(bson.M{"active_at": 1})

	cursor, err := s.TodoCollection.Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, err
	}

	var todos []models.Todo
	err = cursor.All(ctx, &todos)
	if err != nil {
		return nil, ErrNotFound
	}

	return todos, nil
}

func (s *storage) GetOneTodo(ctx context.Context, id string) (todo models.Todo, err error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Todo{}, err
	}
	filter := bson.M{"_id": objId}

	err = s.TodoCollection.FindOne(ctx, filter).Decode(&todo)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Todo{}, ErrNotFound
		}
	}

	return todo, nil
}
