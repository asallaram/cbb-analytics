package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asallaram/cbb-analytics/internal/storage"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Handler struct {
	db *storage.MongoDB
}

func NewHandler(db *storage.MongoDB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetGames(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	date := r.URL.Query().Get("date")
	status := r.URL.Query().Get("status")

	filter := bson.M{}
	if date != "" {
		filter["date"] = date
	}
	if status != "" {
		filter["status"] = status
	}

	cursor, err := h.db.DB.Collection("games").Find(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var games []interface{}
	if err := cursor.All(ctx, &games); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func (h *Handler) GetGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]

	game, err := h.db.GetGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func (h *Handler) GetPlays(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]

	plays, err := h.db.GetPlaysByGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plays)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]

	ctx := r.Context()
	filter := bson.M{"game_id": gameID}

	cursor, err := h.db.DB.Collection("live_stats").Find(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var stats []interface{}
	if err := cursor.All(ctx, &stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetZones(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]

	ctx := r.Context()
	filter := bson.M{"game_id": gameID}

	cursor, err := h.db.DB.Collection("zone_stats").Find(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var zones []interface{}
	if err := cursor.All(ctx, &zones); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}

func (h *Handler) GetInsights(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	insights, err := h.db.GetInsights(gameID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insights)
}
