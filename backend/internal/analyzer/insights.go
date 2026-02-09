package analyzer

import (
	"fmt"
	"strings"
	"time"

	"github.com/asallaram/cbb-analytics/internal/espn"
	"github.com/asallaram/cbb-analytics/internal/models"
)

type InsightGenerator struct {
	plays       []espn.Play
	playerStats map[string]*PlayerStats
	zoneStats   map[string]*ZoneStats
	playerNames map[string]string
}

func NewInsightGenerator(plays []espn.Play) *InsightGenerator {
	return &InsightGenerator{
		plays:       plays,
		playerStats: CalculatePlayerStats(plays),
		zoneStats:   CalculateZoneStats(plays),
		playerNames: extractPlayerNames(plays),
	}
}

func extractPlayerNames(plays []espn.Play) map[string]string {
	names := make(map[string]string)

	for _, play := range plays {
		if len(play.Participants) == 0 {
			continue
		}

		playerID := play.Participants[0].Athlete.ID
		text := play.Text

		// Skip non-player actions
		if strings.Contains(text, "Timeout") ||
			strings.Contains(text, "Team Rebound") ||
			strings.Contains(text, "Official") ||
			strings.Contains(text, "End of") ||
			strings.HasPrefix(text, "Foul on") {
			continue
		}

		// Extract name before action verbs
		words := strings.Fields(text)
		if len(words) < 2 {
			continue
		}

		var nameParts []string
		for _, word := range words {
			lower := strings.ToLower(word)
			// Stop at action words and trailing descriptors
			if lower == "makes" || lower == "misses" || lower == "with" ||
				lower == "turnover" || lower == "bad" || lower == "subbing" ||
				lower == "defensive" || lower == "offensive" || lower == "block." ||
				lower == "steal." || lower == "rebound." {
				break
			}
			nameParts = append(nameParts, word)
		}

		if len(nameParts) >= 2 {
			name := strings.Join(nameParts, " ")
			// Remove trailing periods
			name = strings.TrimRight(name, ".")
			names[playerID] = name
		}
	}

	return names
}

func (ig *InsightGenerator) getPlayerName(playerID string) string {
	if name, exists := ig.playerNames[playerID]; exists {
		return name
	}
	return "Player"
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
		playerName := ig.getPlayerName(playerID)

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
				Title:     fmt.Sprintf("%s on Fire", playerName),
				Message:   fmt.Sprintf("%s shooting %d-%d (%.0f%%) from the field", playerName, stats.FGM, stats.FGA, stats.FGPct),
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
				Title:     fmt.Sprintf("%s Struggling", playerName),
				Message:   fmt.Sprintf("%s shooting %d-%d (%.0f%%) from the field", playerName, stats.FGM, stats.FGA, stats.FGPct),
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
				Title:     fmt.Sprintf("%s Lights Out from Three", playerName),
				Message:   fmt.Sprintf("%s shooting %d-%d (%.0f%%) from beyond the arc", playerName, stats.ThreePM, stats.ThreePA, stats.ThreePct),
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
		playerName := ig.getPlayerName(playerID)

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
					Message:   fmt.Sprintf("%s 0-%d from %s", playerName, data.Attempts, zone),
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
					Message:   fmt.Sprintf("%s %d-%d (%.0f%%) from %s", playerName, data.Makes, data.Attempts, data.Pct, zone),
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
		playerName := ig.getPlayerName(playerID)

		if stats.Turnovers >= 4 {
			insights = append(insights, models.Insight{
				GameID:    gameID,
				Timestamp: time.Now(),
				Type:      "turnover_trouble",
				Category:  "turnovers",
				Severity:  "medium",
				Title:     fmt.Sprintf("%s Turnover Issues", playerName),
				Message:   fmt.Sprintf("%s with %d turnovers", playerName, stats.Turnovers),
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
				Title:     fmt.Sprintf("%s in Foul Trouble", playerName),
				Message:   fmt.Sprintf("%s with %d fouls", playerName, stats.Fouls),
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
