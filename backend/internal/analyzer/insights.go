package analyzer

import (
	"fmt"
	"time"

	"github.com/asallaram/cbb-analytics/internal/espn"
	"github.com/asallaram/cbb-analytics/internal/models"
)

type InsightGenerator struct {
	plays       []espn.Play
	playerStats map[string]*PlayerStats
	zoneStats   map[string]*ZoneStats
}

func NewInsightGenerator(plays []espn.Play) *InsightGenerator {
	return &InsightGenerator{
		plays:       plays,
		playerStats: CalculatePlayerStats(plays),
		zoneStats:   CalculateZoneStats(plays),
	}
}

func (ig *InsightGenerator) GenerateInsights(gameID string) []models.Insight {
	var insights []models.Insight

	insights = append(insights, ig.detectHotColdPlayers(gameID)...)
	insights = append(insights, ig.detectZonePerformance(gameID)...)
	insights = append(insights, ig.detectMomentum(gameID)...)
	insights = append(insights, ig.detectStruggling(gameID)...)

	return insights
}

func (ig *InsightGenerator) detectHotColdPlayers(gameID string) []models.Insight {
	var insights []models.Insight

	for playerID, stats := range ig.playerStats {
		if stats.FGA < 5 {
			continue
		}

		if stats.FGPct >= 60 {
			insights = append(insights, models.Insight{
				GameID:    gameID,
				Timestamp: time.Now(),
				Type:      "player_hot",
				Category:  "shooting",
				Severity:  "high",
				Title:     "Player on Fire",
				Message:   fmt.Sprintf("Player shooting %d-%d (%.0f%%) from the field", stats.FGM, stats.FGA, stats.FGPct),
				Context: models.Context{
					PlayerID: playerID,
					TeamID:   stats.TeamID,
					Stats: map[string]interface{}{
						"fgm":    stats.FGM,
						"fga":    stats.FGA,
						"fg_pct": stats.FGPct,
					},
				},
			})
		}

		if stats.FGPct <= 30 && stats.FGA >= 8 {
			insights = append(insights, models.Insight{
				GameID:    gameID,
				Timestamp: time.Now(),
				Type:      "player_cold",
				Category:  "shooting",
				Severity:  "medium",
				Title:     "Player Struggling",
				Message:   fmt.Sprintf("Player shooting %d-%d (%.0f%%) from the field", stats.FGM, stats.FGA, stats.FGPct),
				Context: models.Context{
					PlayerID: playerID,
					TeamID:   stats.TeamID,
					Stats: map[string]interface{}{
						"fgm":    stats.FGM,
						"fga":    stats.FGA,
						"fg_pct": stats.FGPct,
					},
				},
			})
		}

		if stats.ThreePA >= 4 && stats.ThreePct >= 50 {
			insights = append(insights, models.Insight{
				GameID:    gameID,
				Timestamp: time.Now(),
				Type:      "three_point_hot",
				Category:  "shooting",
				Severity:  "high",
				Title:     "Three-Point Sniper",
				Message:   fmt.Sprintf("Player shooting %d-%d (%.0f%%) from three", stats.ThreePM, stats.ThreePA, stats.ThreePct),
				Context: models.Context{
					PlayerID: playerID,
					TeamID:   stats.TeamID,
					Stats: map[string]interface{}{
						"three_pm":  stats.ThreePM,
						"three_pa":  stats.ThreePA,
						"three_pct": stats.ThreePct,
					},
				},
			})
		}
	}

	return insights
}

func (ig *InsightGenerator) detectZonePerformance(gameID string) []models.Insight {
	var insights []models.Insight

	for playerID, zoneStats := range ig.zoneStats {
		for zone, data := range zoneStats.Zones {
			if data.Attempts < 3 {
				continue
			}

			if data.Pct == 0 {
				insights = append(insights, models.Insight{
					GameID:    gameID,
					Timestamp: time.Now(),
					Type:      "zone_cold",
					Category:  "zone_shooting",
					Severity:  "medium",
					Title:     "Zone Ice Cold",
					Message:   fmt.Sprintf("Player 0-%d from %s", data.Attempts, zone),
					Context: models.Context{
						PlayerID: playerID,
						TeamID:   zoneStats.TeamID,
						Zone:     zone,
						Stats: map[string]interface{}{
							"makes":    data.Makes,
							"attempts": data.Attempts,
							"pct":      data.Pct,
						},
					},
				})
			}

			if data.Pct >= 60 && data.Attempts >= 3 {
				insights = append(insights, models.Insight{
					GameID:    gameID,
					Timestamp: time.Now(),
					Type:      "zone_hot",
					Category:  "zone_shooting",
					Severity:  "high",
					Title:     "Zone Dominant",
					Message:   fmt.Sprintf("Player %d-%d (%.0f%%) from %s", data.Makes, data.Attempts, data.Pct, zone),
					Context: models.Context{
						PlayerID: playerID,
						TeamID:   zoneStats.TeamID,
						Zone:     zone,
						Stats: map[string]interface{}{
							"makes":    data.Makes,
							"attempts": data.Attempts,
							"pct":      data.Pct,
						},
					},
				})
			}
		}
	}

	return insights
}

func (ig *InsightGenerator) detectMomentum(gameID string) []models.Insight {
	var insights []models.Insight

	if len(ig.plays) < 10 {
		return insights
	}

	recentPlays := ig.plays
	if len(ig.plays) > 20 {
		recentPlays = ig.plays[len(ig.plays)-20:]
	}

	teamScores := make(map[string]int)
	for _, play := range recentPlays {
		if play.ScoringPlay && play.Team != nil {
			teamScores[play.Team.ID] += play.ScoreValue
		}
	}

	for teamID, points := range teamScores {
		if points >= 10 {
			insights = append(insights, models.Insight{
				GameID:    gameID,
				Timestamp: time.Now(),
				Type:      "momentum_run",
				Category:  "momentum",
				Severity:  "high",
				Title:     "Scoring Run",
				Message:   fmt.Sprintf("Team on a run with %d points in last 20 plays", points),
				Context: models.Context{
					TeamID: teamID,
					Stats: map[string]interface{}{
						"points":     points,
						"play_count": len(recentPlays),
					},
				},
			})
		}
	}

	return insights
}

func (ig *InsightGenerator) detectStruggling(gameID string) []models.Insight {
	var insights []models.Insight

	for playerID, stats := range ig.playerStats {
		if stats.Turnovers >= 4 {
			insights = append(insights, models.Insight{
				GameID:    gameID,
				Timestamp: time.Now(),
				Type:      "turnover_trouble",
				Category:  "turnovers",
				Severity:  "medium",
				Title:     "Turnover Issues",
				Message:   fmt.Sprintf("Player with %d turnovers", stats.Turnovers),
				Context: models.Context{
					PlayerID: playerID,
					TeamID:   stats.TeamID,
					Stats: map[string]interface{}{
						"turnovers": stats.Turnovers,
					},
				},
			})
		}

		if stats.Fouls >= 4 {
			insights = append(insights, models.Insight{
				GameID:    gameID,
				Timestamp: time.Now(),
				Type:      "foul_trouble",
				Category:  "fouls",
				Severity:  "high",
				Title:     "Foul Trouble",
				Message:   fmt.Sprintf("Player with %d fouls", stats.Fouls),
				Context: models.Context{
					PlayerID: playerID,
					TeamID:   stats.TeamID,
					Stats: map[string]interface{}{
						"fouls": stats.Fouls,
					},
				},
			})
		}
	}

	return insights
}
