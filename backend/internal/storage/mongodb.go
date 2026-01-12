package storage

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	if err := createIndexes(db); err != nil {
		return nil, err
	}

	return &MongoDB{
		client: client,
		DB:     db,
	}, nil
}

func createIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Indexes for games collection
	gamesCollection := db.Collection("games")
	_, err := gamesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "_id", Value: 1}},
	})
	if err != nil {
		return err
	}

	// Indexes for plays collection
	playsCollection := db.Collection("plays")
	_, err = playsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "game_id", Value: 1}},
	})
	if err != nil {
		return err
	}

	_, err = playsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "id", Value: 1}},
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
