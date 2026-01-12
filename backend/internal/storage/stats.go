package storage

import (
	"context"

	"github.com/asallaram/cbb-analytics/internal/analyzer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) UpsertPlayerStats(stats map[string]*analyzer.PlayerStats) error {
	ctx := context.Background()

	for _, stat := range stats {
		filter := bson.M{
			"game_id":   stat.GameID,
			"player_id": stat.PlayerID,
		}
		update := bson.M{"$set": stat}
		opts := options.Update().SetUpsert(true)

		_, err := m.DB.Collection("live_stats").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
func (m *MongoDB) UpsertZoneStats(stats map[string]*analyzer.ZoneStats) error {
	ctx := context.Background()

	for _, stat := range stats {
		filter := bson.M{
			"game_id":   stat.GameID,
			"player_id": stat.PlayerID,
		}
		update := bson.M{"$set": stat}
		opts := options.Update().SetUpsert(true)

		_, err := m.DB.Collection("zone_stats").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
