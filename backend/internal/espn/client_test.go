package espn

import (
	"fmt"
	"testing"
	"time"
)

func TestGetScoreboard(t *testing.T) {
	client := NewClient("https://site.api.espn.com")

	// Get today's games
	today := time.Now().Format("20060102")

	fmt.Printf("Fetching games for date: %s\n", today)

	scoreboard, err := client.GetScoreboard(today)
	if err != nil {
		t.Fatalf("Error fetching scoreboard: %v", err)
	}

	fmt.Printf("\n✅ SUCCESS! Found %d games\n\n", len(scoreboard.Events))

	// Print first 5 games
	for i, event := range scoreboard.Events {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, event.Name)
		fmt.Printf("   ID: %s\n", event.ID)
		fmt.Printf("   Status: %s\n", event.Status.Type.State)
		fmt.Printf("   Time: %s\n\n", event.Date)
	}
}
