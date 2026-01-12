package analyzer

import (
	"strings"

	"github.com/asallaram/cbb-analytics/internal/espn"
)

type PlayerStats struct {
	GameID    string  `bson:"game_id" json:"game_id"`
	PlayerID  string  `bson:"player_id" json:"player_id"`
	TeamID    string  `bson:"team_id" json:"team_id"`
	Points    int     `bson:"points" json:"points"`
	FGM       int     `bson:"fgm" json:"fgm"`
	FGA       int     `bson:"fga" json:"fga"`
	FGPct     float64 `bson:"fg_pct" json:"fg_pct"`
	ThreePM   int     `bson:"three_pm" json:"three_pm"`
	ThreePA   int     `bson:"three_pa" json:"three_pa"`
	ThreePct  float64 `bson:"three_pct" json:"three_pct"`
	FTM       int     `bson:"ftm" json:"ftm"`
	FTA       int     `bson:"fta" json:"fta"`
	FTPct     float64 `bson:"ft_pct" json:"ft_pct"`
	Rebounds  int     `bson:"rebounds" json:"rebounds"`
	Assists   int     `bson:"assists" json:"assists"`
	Steals    int     `bson:"steals" json:"steals"`
	Blocks    int     `bson:"blocks" json:"blocks"`
	Turnovers int     `bson:"turnovers" json:"turnovers"`
	Fouls     int     `bson:"fouls" json:"fouls"`
}

func CalculatePlayerStats(plays []espn.Play) map[string]*PlayerStats {
	stats := make(map[string]*PlayerStats)

	for _, play := range plays {
		if len(play.Participants) == 0 {
			continue
		}

		playerID := play.Participants[0].Athlete.ID

		if _, exists := stats[playerID]; !exists {
			stats[playerID] = &PlayerStats{
				GameID:   play.ID[:9],
				PlayerID: playerID,
				TeamID:   getTeamID(play),
			}
		}

		s := stats[playerID]
		playType := strings.ToLower(play.Type.Text)
		playText := strings.ToLower(play.Text)

		switch {
		case strings.Contains(playType, "jumpshot") || strings.Contains(playType, "layupshot") || strings.Contains(playType, "dunkshot"):
			if strings.Contains(playText, "makes") {
				s.FGM++
				s.Points += play.ScoreValue
				if strings.Contains(playText, "three point") {
					s.ThreePM++
				}
			}
			s.FGA++
			if strings.Contains(playText, "three point") {
				s.ThreePA++
			}

		case strings.Contains(playType, "freethrow"):
			if strings.Contains(playText, "makes") {
				s.FTM++
				s.Points++
			}
			s.FTA++

		case strings.Contains(playType, "rebound"):
			s.Rebounds++

		case strings.Contains(playText, "assists"):
			s.Assists++

		case strings.Contains(playType, "steal"):
			s.Steals++

		case strings.Contains(playType, "block"):
			s.Blocks++

		case strings.Contains(playType, "turnover"):
			s.Turnovers++

		case strings.Contains(playType, "foul"):
			s.Fouls++
		}
	}

	for _, s := range stats {
		if s.FGA > 0 {
			s.FGPct = float64(s.FGM) / float64(s.FGA) * 100
		}
		if s.ThreePA > 0 {
			s.ThreePct = float64(s.ThreePM) / float64(s.ThreePA) * 100
		}
		if s.FTA > 0 {
			s.FTPct = float64(s.FTM) / float64(s.FTA) * 100
		}
	}

	return stats
}

func getTeamID(play espn.Play) string {
	if play.Team != nil {
		return play.Team.ID
	}
	return ""
}
