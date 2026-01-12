package espn

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestGetGameSummary(t *testing.T) {
	client := NewClient("https://site.api.espn.com")

	gameID := "401822893"

	fmt.Printf("Fetching game summary for: %s\n", gameID)

	summary, err := client.GetGameSummary(gameID)
	if err != nil {
		t.Fatalf("Error fetching game summary: %v", err)
	}

	fmt.Printf("\n✅ Game ID: %s\n", summary.Header.ID)

	data, _ := json.MarshalIndent(summary, "", "  ")
	os.WriteFile("/tmp/game_summary.json", data, 0644)
	fmt.Println("✅ Saved full response to /tmp/game_summary.json")

	fmt.Printf("\nChecking for plays...\n")

	t.Logf("Summary structure: %+v", summary)
}
