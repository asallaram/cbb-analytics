# 🏀 Live College Basketball Analytics Platform

Real-time CBB analytics with play-by-play tracking, zone detection, and automated insights. Plan to eventually host this via website but for now just local hosting.

## Tech Stack

- **Backend:** Go, MongoDB, Docker
- **Frontend:** Angular, TypeScript
- **APIs:** ESPN live game data
- **Analytics:** Zone detection, live stats, insight generation

## Features

- Real-time game tracking with 30-second polling
- Play-by-play analysis with court zone detection (7 zones)
- Live player statistics (FG%, 3P%, rebounds, assists, turnovers, fouls)
- Automated insight generation (hot/cold players, zone performance, foul trouble)
- RESTful API with 6 endpoints
- ESPN-style scoreboard UI

## Prerequisites

- Go 1.23+
- Node.js 18+
- Docker & Docker Compose
- MongoDB (or use provided Docker setup)

## Quick Start

### 1. Clone the repository
```bash
git clone https://github.com/asallaram/cbb-analytics.git
cd cbb-analytics
```

### 2. Start MongoDB
```bash
docker-compose up -d
```

### 3. Start Backend

**API Server:**
```bash
cd backend
go run cmd/api/main.go
```
API runs on `http://localhost:8080`

**Poller (separate terminal):**
```bash
cd backend
go run cmd/poller/main.go
```
Polls ESPN every 30 seconds for live games

### 4. Start Frontend
```bash
cd frontend
npm install
ng serve
```
Frontend runs on `http://localhost:4200`

## API Endpoints
```
GET /api/games?status=in          # Get live games
GET /api/games/:id                # Get game details
GET /api/games/:id/plays          # Get play-by-play
GET /api/games/:id/stats          # Get player stats
GET /api/games/:id/zones          # Get zone shooting stats
GET /api/games/:id/insights       # Get automated insights
```

## Project Structure
```
cbb-analytics/
├── backend/
│   ├── cmd/
│   │   ├── api/          # REST API server
│   │   └── poller/       # ESPN data poller
│   ├── internal/
│   │   ├── analyzer/     # Stats & insights engine
│   │   ├── api/          # HTTP handlers
│   │   ├── espn/         # ESPN API client
│   │   ├── models/       # Data models
│   │   └── storage/      # MongoDB operations
│   └── go.mod
├── frontend/
│   └── src/
│       └── app/
│           ├── components/
│           └── services/
└── docker-compose.yml
```

## Future Features

- Historical game browser
- Player comparison tools
- Advanced analytics (momentum tracking, lineup analysis)
- Shot chart visualizations
- Season-long performance tracking

## License

MIT
