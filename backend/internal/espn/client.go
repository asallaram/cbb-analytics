package espn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetScoreboard fetches all games for a given date
func (c *Client) GetScoreboard(date string) (*ScoreboardResponse, error) {
	url := fmt.Sprintf("%s/apis/site/v2/sports/basketball/mens-college-basketball/scoreboard?dates=%s&limit=500",
		c.BaseURL, date)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var scoreboard ScoreboardResponse
	if err := json.Unmarshal(body, &scoreboard); err != nil {
		return nil, err
	}

	return &scoreboard, nil
}

// GetGameSummary fetches complete game data including plays, box score, etc.
func (c *Client) GetGameSummary(gameID string) (*GameSummary, error) {
	url := fmt.Sprintf("%s/apis/site/v2/sports/basketball/mens-college-basketball/summary?event=%s",
		c.BaseURL, gameID)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var summary GameSummary
	if err := json.Unmarshal(body, &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}
