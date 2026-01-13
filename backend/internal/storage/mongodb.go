package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/asallaram/cbb-analytics/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbName)

	if err := createIndexes(db); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &MongoDB{
		client: client,
		DB:     db,
	}, nil
}

func createIndexes(db *mongo.Database) error {
	ctx := context.Background()

	_, err := db.Collection("games").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "date", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("plays").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "game_id", Value: 1}, {Key: "sequence", Value: 1}}},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("insights").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "game_id", Value: 1}, {Key: "timestamp", Value: -1}}},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("live_stats").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "game_id", Value: 1}, {Key: "player_id", Value: 1}}},
	})

	return err
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

func (m *MongoDB) UpsertGame(game *models.Game) error {
	ctx := context.Background()

	filter := bson.M{"_id": game.ID}
	update := bson.M{"$set": game}
	opts := options.Update().SetUpsert(true)

	_, err := m.DB.Collection("games").UpdateOne(ctx, filter, update, opts)
	return err
}

func (m *MongoDB) GetGame(gameID string) (*models.Game, error) {
	ctx := context.Background()

	var game models.Game
	err := m.DB.Collection("games").FindOne(ctx, bson.M{"_id": gameID}).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}
