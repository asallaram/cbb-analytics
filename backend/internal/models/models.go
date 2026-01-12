package models

import (
	"time"
)

type Game struct {
	ID            string    `bson:"id" json:"id"`
	Date          string    `bson:"date" json:"date"`
	HomeTeamID    string    `bson:"home_team_id" json:"home_team_id"`
	HomeTeamName  string    `bson:"home_team_name" json:"home_team_name"`
	AwayTeamID    string    `bson:"away_team_id" json:"away_team_id"`
	AwayTeamName  string    `bson:"away_team_name" json:"away_team_name"`
	Status        string    `bson:"status" json:"status"`
	CurrentPeriod int       `bson:"current_period" json:"current_period"`
	CurrentClock  string    `bson:"current_clock" json:"current_clock"`
	HomeScore     int       `bson:"home_score" json:"home_score"`
	AwayScore     int       `bson:"away_score" json:"away_score"`
	LastUpdated   time.Time `bson:"last_updated" json:"last_updated"`
}

type Score struct {
	Home int `bson:"home" json:"home"`
	Away int `bson:"away" json:"away"`
}

type Team struct {
	ID           string `bson:"_id" json:"id"`
	Name         string `bson:"name" json:"name"`
	DisplayName  string `bson:"display_name" json:"display_name"`
	Abbreviation string `bson:"abbreviation" json:"abbreviation"`
	Logo         string `bson:"logo" json:"logo"`
}

type Insight struct {
	GameID    string    `bson:"game_id" json:"game_id"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	GameClock string    `bson:"game_clock" json:"game_clock"`
	Type      string    `bson:"type" json:"type"`
	Category  string    `bson:"category" json:"category"`
	Severity  string    `bson:"severity" json:"severity"`
	Title     string    `bson:"title" json:"title"`
	Message   string    `bson:"message" json:"message"`
	Context   Context   `bson:"context" json:"context"`
}

type Context struct {
	TeamID   string                 `bson:"team_id,omitempty" json:"team_id,omitempty"`
	PlayerID string                 `bson:"player_id,omitempty" json:"player_id,omitempty"`
	Zone     string                 `bson:"zone,omitempty" json:"zone,omitempty"`
	Stats    map[string]interface{} `bson:"stats,omitempty" json:"stats,omitempty"`
}
type Play struct {
	ID             string   `bson:"id" json:"id"`
	GameID         string   `bson:"game_id" json:"game_id"`
	SequenceNumber string   `bson:"sequence_number" json:"sequence_number"`
	Type           string   `bson:"type" json:"type"`
	TypeID         string   `bson:"type_id" json:"type_id"`
	Text           string   `bson:"text" json:"text"`
	Period         int      `bson:"period" json:"period"`
	Clock          string   `bson:"clock" json:"clock"`
	AwayScore      int      `bson:"away_score" json:"away_score"`
	HomeScore      int      `bson:"home_score" json:"home_score"`
	ScoringPlay    bool     `bson:"scoring_play" json:"scoring_play"`
	ScoreValue     int      `bson:"score_value" json:"score_value"`
	ShootingPlay   bool     `bson:"shooting_play" json:"shooting_play"`
	CoordinateX    *float64 `bson:"coordinate_x,omitempty" json:"coordinate_x,omitempty"`
	CoordinateY    *float64 `bson:"coordinate_y,omitempty" json:"coordinate_y,omitempty"`
	TeamID         string   `bson:"team_id,omitempty" json:"team_id,omitempty"`
	PlayerIDs      []string `bson:"player_ids,omitempty" json:"player_ids,omitempty"`
	Timestamp      string   `bson:"timestamp" json:"timestamp"`
}
