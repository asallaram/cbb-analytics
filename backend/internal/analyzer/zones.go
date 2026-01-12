package analyzer

import (
	"github.com/asallaram/cbb-analytics/internal/espn"
)

type ZoneStats struct {
	GameID   string              `bson:"game_id" json:"game_id"`
	TeamID   string              `bson:"team_id" json:"team_id"`
	PlayerID string              `bson:"player_id,omitempty" json:"player_id,omitempty"`
	Zones    map[string]ZoneData `bson:"zones" json:"zones"`
}

type ZoneData struct {
	Makes    int     `bson:"makes" json:"makes"`
	Attempts int     `bson:"attempts" json:"attempts"`
	Pct      float64 `bson:"pct" json:"pct"`
}

func GetZone(x, y float64) string {
	if x < 10 && y < 10 {
		return "left_corner_3"
	}
	if x > 40 && y < 10 {
		return "right_corner_3"
	}
	if x < 20 && y > 10 && y < 24 {
		return "left_wing_3"
	}
	if x > 30 && y > 10 && y < 24 {
		return "right_wing_3"
	}
	if x >= 20 && x <= 30 && y > 20 {
		return "top_key_3"
	}
	if x >= 15 && x <= 35 && y < 10 {
		return "paint"
	}
	return "mid_range"
}

func CalculateZoneStats(plays []espn.Play) map[string]*ZoneStats {
	stats := make(map[string]*ZoneStats)

	for _, play := range plays {
		if !play.ShootingPlay || play.Coordinate == nil {
			continue
		}

		if len(play.Participants) == 0 {
			continue
		}

		playerID := play.Participants[0].Athlete.ID
		teamID := getTeamIDFromPlay(play)

		key := playerID
		if _, exists := stats[key]; !exists {
			stats[key] = &ZoneStats{
				GameID:   play.ID[:9],
				TeamID:   teamID,
				PlayerID: playerID,
				Zones:    make(map[string]ZoneData),
			}
		}

		zone := GetZone(play.Coordinate.X, play.Coordinate.Y)
		zoneData := stats[key].Zones[zone]
		zoneData.Attempts++
		if play.ScoringPlay {
			zoneData.Makes++
		}
		if zoneData.Attempts > 0 {
			zoneData.Pct = float64(zoneData.Makes) / float64(zoneData.Attempts) * 100
		}
		stats[key].Zones[zone] = zoneData
	}

	return stats
}

func getTeamIDFromPlay(play espn.Play) string {
	if play.Team != nil {
		return play.Team.ID
	}
	return ""
}
