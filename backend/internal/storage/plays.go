package storage

import (
	"context"

	"github.com/asallaram/cbb-analytics/internal/espn"
	"github.com/asallaram/cbb-analytics/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) UpsertPlays(gameID string, plays []espn.Play) error {
	ctx := context.Background()

	if len(plays) == 0 {
		return nil
	}

	var modelPlays []interface{}
	for _, p := range plays {
		play := models.Play{
			ID:             p.ID,
			GameID:         gameID,
			SequenceNumber: p.SequenceNumber,
			Type:           p.Type.Text,
			TypeID:         p.Type.ID,
			Text:           p.Text,
			Period:         p.Period.Number,
			Clock:          p.Clock.DisplayValue,
			AwayScore:      p.AwayScore,
			HomeScore:      p.HomeScore,
			ScoringPlay:    p.ScoringPlay,
			ScoreValue:     p.ScoreValue,
			ShootingPlay:   p.ShootingPlay,
			Timestamp:      p.Wallclock,
		}

		if p.Coordinate != nil {
			play.CoordinateX = &p.Coordinate.X
			play.CoordinateY = &p.Coordinate.Y
		}

		if p.Team != nil {
			play.TeamID = p.Team.ID
		}

		if len(p.Participants) > 0 {
			for _, participant := range p.Participants {
				play.PlayerIDs = append(play.PlayerIDs, participant.Athlete.ID)
			}
		}

		modelPlays = append(modelPlays, play)
	}

	for _, play := range modelPlays {
		filter := bson.M{"id": play.(models.Play).ID}
		update := bson.M{"$set": play}
		opts := options.Update().SetUpsert(true)

		_, err := m.DB.Collection("plays").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MongoDB) GetPlaysByGame(gameID string) ([]models.Play, error) {
	ctx := context.Background()

	filter := bson.M{"game_id": gameID}
	opts := options.Find().SetSort(bson.M{"sequence_number": 1})

	cursor, err := m.DB.Collection("plays").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plays []models.Play
	if err := cursor.All(ctx, &plays); err != nil {
		return nil, err
	}

	return plays, nil
}
