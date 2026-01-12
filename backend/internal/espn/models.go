package espn

type ScoreboardResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID           string        `json:"id"`
	Date         string        `json:"date"`
	Name         string        `json:"name"`
	Status       Status        `json:"status"`
	Competitions []Competition `json:"competitions"`
}

type Status struct {
	Type struct {
		State       string `json:"state"`
		Description string `json:"description"`
	} `json:"type"`
	Period       int    `json:"period"`
	DisplayClock string `json:"displayClock"`
}

type Competition struct {
	ID          string       `json:"id"`
	Competitors []Competitor `json:"competitors"`
}

type Competitor struct {
	ID       string `json:"id"`
	Team     Team   `json:"team"`
	HomeAway string `json:"homeAway"`
	Score    string `json:"score"`
}

type Team struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
	Logo         string `json:"logo"`
}

// GameSummary with Plays
type GameSummary struct {
	Header   Header   `json:"header"`
	BoxScore BoxScore `json:"boxscore"`
	Plays    []Play   `json:"plays"`
}

type Header struct {
	ID string `json:"id"`
}

type BoxScore struct {
	Teams   []TeamStats   `json:"teams"`
	Players []PlayerStats `json:"players"`
}

type TeamStats struct {
	Team       Team        `json:"team"`
	Statistics []Statistic `json:"statistics"`
}

type PlayerStats struct {
	Team       Team `json:"team"`
	Statistics []struct {
		Athletes []Athlete `json:"athletes"`
	} `json:"statistics"`
}

type Athlete struct {
	Athlete struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"athlete"`
	Stats []string `json:"stats"`
}

type Statistic struct {
	Name         string `json:"name"`
	DisplayValue string `json:"displayValue"`
}

// Play-by-play structures
type Play struct {
	ID             string        `json:"id"`
	SequenceNumber string        `json:"sequenceNumber"`
	Type           PlayType      `json:"type"`
	Text           string        `json:"text"`
	AwayScore      int           `json:"awayScore"`
	HomeScore      int           `json:"homeScore"`
	Period         Period        `json:"period"`
	Clock          Clock         `json:"clock"`
	ScoringPlay    bool          `json:"scoringPlay"`
	ScoreValue     int           `json:"scoreValue"`
	ShootingPlay   bool          `json:"shootingPlay"`
	Coordinate     *Coordinate   `json:"coordinate,omitempty"`
	Team           *Team         `json:"team,omitempty"`
	Participants   []Participant `json:"participants,omitempty"`
	Wallclock      string        `json:"wallclock"`
}

type PlayType struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Period struct {
	Number       int    `json:"number"`
	DisplayValue string `json:"displayValue"`
}

type Clock struct {
	DisplayValue string `json:"displayValue"`
}

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Participant struct {
	Athlete struct {
		ID string `json:"id"`
	} `json:"athlete"`
}
