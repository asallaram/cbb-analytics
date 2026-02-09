package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/asallaram/cbb-analytics/internal/analyzer"
	"github.com/asallaram/cbb-analytics/internal/espn"
	"github.com/asallaram/cbb-analytics/internal/models"
	"github.com/asallaram/cbb-analytics/internal/storage"
)

func main() {
	client := espn.NewClient("https://site.api.espn.com")

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongo, err := storage.NewMongoDB(mongoURI, "cbb_analytics")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongo.Close()

	fmt.Println("🏀 Live Game Poller Started!")
	fmt.Println("Polling ESPN every 30 seconds for live games...")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	pollGames(client, mongo)

	for range ticker.C {
		pollGames(client, mongo)
	}
}

func pollGames(client *espn.Client, mongo *storage.MongoDB) {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	for _, date := range []time.Time{yesterday, today} {
		dateStr := date.Format("20060102")

		scoreboard, err := client.GetScoreboard(dateStr)
		if err != nil {
			log.Printf("Error fetching scoreboard for %s: %v", dateStr, err)
			continue
		}

		if date.Format("20060102") == today.Format("20060102") {
			fmt.Printf("\n[%s] Found %d games today\n", time.Now().Format("15:04:05"), len(scoreboard.Events))
		}

		for _, event := range scoreboard.Events {
			if len(event.Competitions) == 0 {
				continue
			}

			comp := event.Competitions[0]
			status := event.Status.Type.State

			game := convertEventToGame(event, comp)

			if err := mongo.UpsertGame(&game); err != nil {
				log.Printf("Error saving game %s: %v", game.ID, err)
				continue
			}

			if status == "in" && date.Format("20060102") == today.Format("20060102") {
				fmt.Printf("🔴 LIVE: %s (fetching plays...)\n", event.Name)

				summary, err := client.GetGameSummary(game.ID)
				if err != nil {
					log.Printf("Error fetching summary for %s: %v", game.ID, err)
					continue
				}

				if err := mongo.UpsertPlays(game.ID, summary.Plays); err != nil {
					log.Printf("Error saving plays for %s: %v", game.ID, err)
					continue
				}

				stats := analyzer.CalculatePlayerStats(summary.Plays)
				if err := mongo.UpsertPlayerStats(stats); err != nil {
					log.Printf("Error saving stats for %s: %v", game.ID, err)
					continue
				}

				zones := analyzer.CalculateZoneStats(summary.Plays)
				if err := mongo.UpsertZoneStats(zones); err != nil {
					log.Printf("Error saving zones for %s: %v", game.ID, err)
					continue
				}

				generator := analyzer.NewInsightGenerator(summary.Plays)
				insights := generator.GenerateInsights(game.ID)
				if err := mongo.SaveInsights(insights); err != nil {
					log.Printf("Error saving insights for %s: %v", game.ID, err)
					continue
				}

				fmt.Printf("   └─ Saved %d plays, %d players, %d zones, %d insights\n",
					len(summary.Plays), len(stats), len(zones), len(insights))
			}
		}
	}
}

func convertEventToGame(event espn.Event, comp espn.Competition) models.Game {
	game := models.Game{
		ID:            event.ID,
		Date:          event.Date,
		Status:        event.Status.Type.State,
		CurrentPeriod: event.Status.Period,
		CurrentClock:  event.Status.DisplayClock,
		LastUpdated:   time.Now(),
	}

	for _, competitor := range comp.Competitors {
		if competitor.HomeAway == "home" {
			game.HomeTeamID = competitor.Team.ID
			game.HomeTeamName = competitor.Team.DisplayName
			game.HomeScore = parseInt(competitor.Score)
		} else {
			game.AwayTeamID = competitor.Team.ID
			game.AwayTeamName = competitor.Team.DisplayName
			game.AwayScore = parseInt(competitor.Score)
		}
	}

	return game
}

func parseInt(s string) int {
	var val int
	fmt.Sscanf(s, "%d", &val)
	return val
}
