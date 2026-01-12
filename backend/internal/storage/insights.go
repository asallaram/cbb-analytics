package storage

import (
	"context"

	"github.com/asallaram/cbb-analytics/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) SaveInsights(insights []models.Insight) error {
	ctx := context.Background()

	if len(insights) == 0 {
		return nil
	}

	var docs []interface{}
	for _, insight := range insights {
		docs = append(docs, insight)
	}

	_, err := m.DB.Collection("insights").InsertMany(ctx, docs)
	return err
}

func (m *MongoDB) GetInsights(gameID string, limit int) ([]models.Insight, error) {
	ctx := context.Background()

	filter := bson.M{"game_id": gameID}
	opts := options.Find().
		SetSort(bson.M{"timestamp": -1}).
		SetLimit(int64(limit))

	cursor, err := m.DB.Collection("insights").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var insights []models.Insight
	if err := cursor.All(ctx, &insights); err != nil {
		return nil, err
	}

	return insights, nil
}
